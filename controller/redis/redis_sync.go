package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"monitor-trade/model"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
)

// LoadMonitorPairsFromRedis 启动时从Redis加载所有监控数据到本地
func (r *RedisController) LoadMonitorPairsFromRedis() error {
	r.mutexMonitorPairs.Lock()
	defer r.mutexMonitorPairs.Unlock()

	ctx := context.Background()
	pattern := fmt.Sprintf("%s:*", MonitorKey)
	keys, err := r.Client.Keys(ctx, pattern).Result()
	if err != nil {
		return fmt.Errorf("获取监控键列表失败: %v", err)
	}

	log.Printf("从Redis加载监控数据，找到 %d 个键", len(keys))

	for _, key := range keys {
		val, err := r.Client.Get(ctx, key).Result()
		if err != nil {
			log.Printf("获取Redis键 %s 的值失败: %v", key, err)
			continue
		}

		var pairData model.PairMonitorData
		if err := json.Unmarshal([]byte(val), &pairData); err != nil {
			log.Printf("解析Redis键 %s 的值失败: %v", key, err)
			continue
		}

		// 生成本地map的key: pair:direct
		localKey := fmt.Sprintf("%s:%s", pairData.Pair, pairData.Direct)
		r.MonitorPairs[localKey] = pairData
		log.Printf("加载监控数据: %s", localKey)
	}

	log.Printf("监控数据加载完成，本地缓存 %d 条记录", len(r.MonitorPairs))
	return nil
}

// syncMonitorPairFromRedis 从Redis同步单个监控数据到本地
func (r *RedisController) syncMonitorPairFromRedis(redisKey string) {
	r.mutexMonitorPairs.Lock()
	defer r.mutexMonitorPairs.Unlock()

	ctx := context.Background()
	val, err := r.Client.Get(ctx, redisKey).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			// 键不存在，从本地删除
			r.removeMonitorPairByRedisKey(redisKey)
			return
		}
		log.Printf("同步Redis键 %s 失败: %v", redisKey, err)
		return
	}

	var pairData model.PairMonitorData
	if err := json.Unmarshal([]byte(val), &pairData); err != nil {
		log.Printf("解析Redis键 %s 的值失败: %v", redisKey, err)
		return
	}

	// 生成本地map的key: pair:direct
	localKey := fmt.Sprintf("%s:%s", pairData.Pair, pairData.Direct)
	r.MonitorPairs[localKey] = pairData
	log.Printf("同步监控数据: %s", localKey)
}

// removeMonitorPairByRedisKey 根据Redis键从本地删除监控数据
func (r *RedisController) removeMonitorPairByRedisKey(redisKey string) {
	// 解析Redis键格式: monitor:pair:direct
	parts := strings.Split(redisKey, ":")
	if len(parts) >= 3 {
		pair := strings.Join(parts[1:len(parts)-1], ":")
		direct := parts[len(parts)-1]
		localKey := fmt.Sprintf("%s:%s", pair, direct)

		if _, exists := r.MonitorPairs[localKey]; exists {
			delete(r.MonitorPairs, localKey)
			log.Printf("删除本地监控数据: %s", localKey)
		}
	}
}

// EnableRedisKeyspaceNotifications 启用Redis keyspace通知
func (r *RedisController) EnableRedisKeyspaceNotifications() error {
	ctx := context.Background()

	// 检查当前配置
	configs, err := r.Client.ConfigGet(ctx, "notify-keyspace-events").Result()
	if err != nil {
		return fmt.Errorf("获取Redis配置失败: %v", err)
	}

	// 检查是否已经启用了keyspace事件
	if len(configs) >= 2 {
		currentConfig := configs[1].(string)
		if strings.Contains(currentConfig, "K") && strings.Contains(currentConfig, "E") {
			log.Printf("Redis keyspace事件已启用: %s", currentConfig)
			return nil
		}
	}

	// 启用keyspace事件: K = keyspace events, E = keyevent events
	err = r.Client.ConfigSet(ctx, "notify-keyspace-events", "KEA").Err()
	if err != nil {
		return fmt.Errorf("启用Redis keyspace事件失败: %v", err)
	}

	log.Println("已启用Redis keyspace事件通知")
	return nil
}

// StartRedisSync 启动Redis keyspace事件监听，自动同步本地MonitorPairs
func (r *RedisController) StartRedisSync() {
	ctx := r.Client.Context()

	// 订阅keyspace事件，监听所有monitor:*键的变化
	pattern := fmt.Sprintf("__keyspace@0__:%s:*", MonitorKey)
	pubsub := r.Client.PSubscribe(ctx, pattern)
	defer pubsub.Close()

	log.Printf("开始监听Redis keyspace事件: %s", pattern)

	for {
		msg, err := pubsub.ReceiveMessage(ctx)
		if err != nil {
			log.Printf("接收Redis keyspace消息失败: %v", err)
			time.Sleep(time.Second)
			continue
		}

		// 解析Redis键
		pairKey := msg.Channel
		trimStr := "__keyspace@0__:"
		pairKey = strings.TrimPrefix(pairKey, trimStr)

		log.Printf("收到Redis事件: 键=%s, 操作=%s", pairKey, msg.Payload)

		// 处理不同的Redis事件
		switch msg.Payload {
		case "set":
			// 键被设置，同步到本地
			go r.syncMonitorPairFromRedis(pairKey)

		case "expired", "del":
			// 键被删除或过期，从本地删除
			go func() {
				r.mutexMonitorPairs.Lock()
				r.removeMonitorPairByRedisKey(pairKey)
				r.mutexMonitorPairs.Unlock()
			}()

		default:
			// 其他事件暂时忽略
			continue
		}
	}
}
