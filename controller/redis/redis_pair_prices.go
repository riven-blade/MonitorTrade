package redis

import (
	"monitor-trade/model"
)

// UpdatePairPrice 更新交易对价格数据
func (r *RedisController) UpdatePairPrice(pair string, pairData *model.PairData) {
	r.mutexPairPrices.Lock()
	defer r.mutexPairPrices.Unlock()

	if r.PairPrices == nil {
		r.PairPrices = make(map[string]*model.PairData)
	}
	r.PairPrices[pair] = pairData
}

// GetPairPrice 获取交易对价格数据
func (r *RedisController) GetPairPrice(pair string) model.PairData {
	r.mutexPairPrices.RLock()
	defer r.mutexPairPrices.RUnlock()

	if r.PairPrices == nil {
		return model.PairData{}
	}

	if data, exists := r.PairPrices[pair]; exists {
		return *data
	}
	return model.PairData{}
}

// GetAllPairPricesData 获取所有交易对价格数据
func (r *RedisController) GetAllPairPricesData() ([]model.PairData, error) {
	r.mutexPairPrices.RLock()
	defer r.mutexPairPrices.RUnlock()

	if r.PairPrices == nil {
		return []model.PairData{}, nil
	}

	var pairsData []model.PairData
	for _, data := range r.PairPrices {
		if r.IsWatchedPair(data.Pair) {
			pairsData = append(pairsData, *data)
		}
	}
	return pairsData, nil
}
