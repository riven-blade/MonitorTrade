package model

// PairData 定义了从 Redis 获取的数据结构
type PairData struct {
	Timestamp string  `json:"timestamp"`
	Pair      string  `json:"pair"`
	BidPrice  float64 `json:"bid_price"`
	AskPrice  float64 `json:"ask_price"`
	Close     float64 `json:"close"`
}

type PairMonitorData struct {
	Timestamp string  `json:"timestamp"`
	Pair      string  `json:"pair"`
	Direct    string  `json:"direct"`
	Price     float64 `json:"price"`
}

type PairMonitorDataDetail struct {
	PairData        PairData        `json:"pair_data"`
	PairMonitorData PairMonitorData `json:"pair_monitor_data"`
	TTL             float64         `json:"ttl"` // TTL in seconds
}

type PairMonitorDataWithTTL struct {
	PairMonitorData PairMonitorData
	TTL             float64 `json:"ttl"` // TTL in seconds
}
