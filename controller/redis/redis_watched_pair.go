package redis

// SetWatchedPairs 设置需要监听的交易对
func (r *RedisController) SetWatchedPairs(pairs []string) {
	r.mutexWatchedPairs.Lock()
	defer r.mutexWatchedPairs.Unlock()

	r.WatchedPairs = pairs
}

// GetWatchedPairs 获取需要监听的交易对
func (r *RedisController) GetWatchedPairs() []string {
	r.mutexWatchedPairs.RLock()
	defer r.mutexWatchedPairs.RUnlock()

	return r.WatchedPairs
}

// IsWatchedPair 检查是否为监听的交易对
func (r *RedisController) IsWatchedPair(pair string) bool {
	r.mutexWatchedPairs.RLock()
	defer r.mutexWatchedPairs.RUnlock()

	if len(r.WatchedPairs) == 0 {
		return true // 如果没有设置监听列表，则监听所有交易对
	}

	for _, watchedPair := range r.WatchedPairs {
		if pair == watchedPair {
			return true
		}
	}
	return false
}
