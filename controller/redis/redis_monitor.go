package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"monitor-trade/model"
	"time"
)

// ===== 本地 MonitorPairs 操作（主要数据源） =====

// GetMonitorPair 获取本地监控数据
func (r *RedisController) GetMonitorPair(pair, direct string) (model.PairMonitorData, bool) {
	r.mutexMonitorPairs.RLock()
	defer r.mutexMonitorPairs.RUnlock()

	localKey := fmt.Sprintf("%s:%s", pair, direct)
	data, exists := r.MonitorPairs[localKey]
	return data, exists
}

// SetMonitorPair 设置本地监控数据并同步到Redis
func (r *RedisController) SetMonitorPair(data model.PairMonitorData, direct string) error {
	// 先更新本地数据
	r.mutexMonitorPairs.Lock()
	localKey := fmt.Sprintf("%s:%s", data.Pair, direct)
	data.Direct = direct
	data.Timestamp = time.Now().Format("2006-01-02 15:04:05")
	r.MonitorPairs[localKey] = data
	r.mutexMonitorPairs.Unlock()

	// 同步到Redis
	if err := r.SetPairDataToRedis(data, direct); err != nil {
		// Redis 写入失败，但本地数据已更新，记录错误但不回滚
		log.Printf("同步监控数据到Redis失败 %s: %v", localKey, err)
		return err
	}

	log.Printf("设置监控数据: %s", localKey)
	return nil
}

// DeleteMonitorPair 删除本地监控数据并同步到Redis
func (r *RedisController) DeleteMonitorPair(pair string, direct string) {
	// 先删除本地数据
	r.mutexMonitorPairs.Lock()
	localKey := fmt.Sprintf("%s:%s", pair, direct)
	if _, exists := r.MonitorPairs[localKey]; exists {
		delete(r.MonitorPairs, localKey)
		log.Printf("删除监控数据: %s", localKey)
	}
	r.mutexMonitorPairs.Unlock()

	// 同步删除Redis
	r.deletePairDataRedis(pair, direct)
}

// HasMonitorPair 检查是否存在监控数据
func (r *RedisController) HasMonitorPair(pair, direct string) bool {
	r.mutexMonitorPairs.RLock()
	defer r.mutexMonitorPairs.RUnlock()

	localKey := fmt.Sprintf("%s:%s", pair, direct)
	_, exists := r.MonitorPairs[localKey]
	return exists
}

// getAllMonitorPairsData 获取所有监控中的交易对数据
func (r *RedisController) GetAllMonitorPairsData(direct string) []model.PairMonitorDataWithTTL {
	ctx := context.Background()

	r.mutexMonitorPairs.RLock()
	localData := make(map[string]model.PairMonitorData)
	for k, v := range r.MonitorPairs {
		localData[k] = v
	}
	r.mutexMonitorPairs.RUnlock()

	var pairsData []model.PairMonitorDataWithTTL
	for localKey, pairData := range localData {
		if direct != "" && pairData.Direct != direct {
			continue
		}

		// 构造对应的Redis key
		redisKey := fmt.Sprintf("%s:%s", MonitorKey, localKey)

		// 从Redis查询TTL信息
		ttl, err := r.Client.TTL(ctx, redisKey).Result()
		if err != nil {
			log.Printf("获取 Redis 键 %s 的TTL失败: %v", redisKey, err)
			// TTL查询失败，设置默认值但不跳过数据
			ttl = -1 * time.Second
		}

		pairsData = append(pairsData, model.PairMonitorDataWithTTL{
			PairMonitorData: pairData,
			TTL:             ttl.Seconds(),
		})
	}

	return pairsData
}

// -------------------  redis 操作的方法 -------------------

// deletePairDataRedis 从 Redis 删除单个交易对数据（内部方法）
func (r *RedisController) deletePairDataRedis(pair string, direct string) {
	ctx := context.Background()
	key := fmt.Sprintf("%s:%s:%s", MonitorKey, pair, direct)
	r.Client.Del(ctx, key).Result()
}

// setPairDataToRedis Redis 设置单个交易对数据（内部方法）
func (r *RedisController) SetPairDataToRedis(data model.PairMonitorData, direct string) error {
	ctx := context.Background()
	key := fmt.Sprintf("%s:%s:%s", MonitorKey, data.Pair, direct)

	dataBytes, err := json.Marshal(&data)
	if err != nil {
		return err
	}

	// 生成随机过期时间：r.conf.Redis.KeyExpire/2 到 r.conf.Redis.KeyExpire 之间
	minExpire := r.conf.Redis.KeyExpire / 2
	maxExpire := r.conf.Redis.KeyExpire
	randomExpire := minExpire + rand.Intn(maxExpire-minExpire+1)

	// 只写入Redis，不更新本地缓存（本地数据是主数据源）
	_, err = r.Client.Set(ctx, key, string(dataBytes), time.Duration(randomExpire)*time.Second).Result()
	return err
}
