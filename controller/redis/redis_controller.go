package redis

import (
	"context"
	"monitor-trade/config"
	"monitor-trade/model"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

const MonitorKey = "monitor"
const TradeKey = "trade" // Redis锁的前缀

type RedisController struct {
	Client            *redis.Client
	conf              *config.Config
	MonitorPairs      map[string]model.PairMonitorData
	PairPrices        map[string]*model.PairData
	WatchedPairs      []string     // 需要监听的交易对
	mutexPairPrices   sync.RWMutex // 保护 PairPrices 的读写锁
	mutexWatchedPairs sync.RWMutex // 保护 WatchedPairs 的读写锁
	mutexMonitorPairs sync.RWMutex // 保护 MonitorPairs 的读写锁
}

func NewRedisController(conf *config.Config) *RedisController {
	rdb := redis.NewClient(&redis.Options{
		Addr:     conf.Redis.Addr,
		Password: conf.Redis.Password,
		DB:       conf.Redis.DB,
	})

	return &RedisController{
		Client:            rdb,
		conf:              conf,
		MonitorPairs:      make(map[string]model.PairMonitorData, 1000),
		PairPrices:        make(map[string]*model.PairData, 1000),
		mutexWatchedPairs: sync.RWMutex{},
		mutexPairPrices:   sync.RWMutex{},
		mutexMonitorPairs: sync.RWMutex{},
	}
}

// AcquireTradeLock 获取交易锁
func (r *RedisController) AcquireTradeLock(pair string) bool {
	ctx := context.Background()
	key := TradeKey + ":" + pair

	// 尝试设置锁，TTL为120秒
	result, err := r.Client.SetNX(ctx, key, "locked", 120*time.Second).Result()
	if err != nil {
		return false
	}

	return result
}

// ReleaseTradeLock 释放交易锁
func (r *RedisController) ReleaseTradeLock(pair string) {
	ctx := context.Background()
	key := TradeKey + ":" + pair
	r.Client.Del(ctx, key)
}
