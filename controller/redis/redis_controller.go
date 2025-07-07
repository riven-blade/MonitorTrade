package redis

import (
	"monitor-trade/config"
	"monitor-trade/model"
	"sync"

	"github.com/go-redis/redis/v8"
)

const MonitorKey = "monitor"

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
