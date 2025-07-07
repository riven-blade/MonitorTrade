package model

type TradePosition struct {
	TradeId              int          `json:"trade_id"`
	Pair                 string       `json:"pair"`
	BaseCurrency         string       `json:"base_currency"`
	QuoteCurrency        string       `json:"quote_currency"`
	IsOpen               bool         `json:"is_open"`
	IsShort              bool         `json:"is_short"`
	Exchange             string       `json:"exchange"`
	Amount               float64      `json:"amount"`
	AmountRequested      float64      `json:"amount_requested"`
	StakeAmount          float64      `json:"stake_amount"`
	MaxStakeAmount       float64      `json:"max_stake_amount"`
	Strategy             string       `json:"strategy"`
	EnterTag             string       `json:"enter_tag"`
	Timeframe            int          `json:"timeframe"`
	FeeOpen              float64      `json:"fee_open"`
	FeeOpenCost          float64      `json:"fee_open_cost"`
	FeeOpenCurrency      string       `json:"fee_open_currency"`
	FeeClose             float64      `json:"fee_close"`
	FeeCloseCost         *float64     `json:"fee_close_cost"`
	FeeCloseCurrency     *string      `json:"fee_close_currency"`
	OpenDate             string       `json:"open_date"`
	OpenTimestamp        int64        `json:"open_timestamp"`
	OpenFillDate         string       `json:"open_fill_date"`
	OpenFillTimestamp    int64        `json:"open_fill_timestamp"`
	OpenRate             float64      `json:"open_rate"`
	OpenRateRequested    float64      `json:"open_rate_requested"`
	OpenTradeValue       float64      `json:"open_trade_value"`
	CloseDate            *string      `json:"close_date"`
	CloseTimestamp       *int64       `json:"close_timestamp"`
	CloseRate            *float64     `json:"close_rate"`
	CloseRateRequested   *float64     `json:"close_rate_requested"`
	CloseProfit          *float64     `json:"close_profit"`
	CloseProfitPct       *float64     `json:"close_profit_pct"`
	CloseProfitAbs       *float64     `json:"close_profit_abs"`
	ProfitRatio          float64      `json:"profit_ratio"`
	ProfitPct            float64      `json:"profit_pct"`
	ProfitAbs            float64      `json:"profit_abs"`
	ProfitFiat           float64      `json:"profit_fiat"`
	RealizedProfit       float64      `json:"realized_profit"`
	RealizedProfitRatio  *float64     `json:"realized_profit_ratio"`
	ExitReason           *string      `json:"exit_reason"`
	ExitOrderStatus      *string      `json:"exit_order_status"`
	StopLossAbs          float64      `json:"stop_loss_abs"`
	StopLossRatio        float64      `json:"stop_loss_ratio"`
	StopLossPct          float64      `json:"stop_loss_pct"`
	StoplossLastUpdate   *string      `json:"stoploss_last_update"`
	StoplossLastUpdateTs *int64       `json:"stoploss_last_update_timestamp"`
	InitialStopLossAbs   float64      `json:"initial_stop_loss_abs"`
	InitialStopLossRatio float64      `json:"initial_stop_loss_ratio"`
	InitialStopLossPct   float64      `json:"initial_stop_loss_pct"`
	MinRate              float64      `json:"min_rate"`
	MaxRate              float64      `json:"max_rate"`
	HasOpenOrders        bool         `json:"has_open_orders"`
	Orders               []TradeOrder `json:"orders"`
	Leverage             float64      `json:"leverage"`
	InterestRate         float64      `json:"interest_rate"`
	LiquidationPrice     float64      `json:"liquidation_price"`
	FundingFees          float64      `json:"funding_fees"`
	TradingMode          string       `json:"trading_mode"`
	AmountPrecision      float64      `json:"amount_precision"`
	PricePrecision       float64      `json:"price_precision"`
	PrecisionMode        int          `json:"precision_mode"`
	StoplossCurrentDist  float64      `json:"stoploss_current_dist"`
	StoplossCurrentPct   float64      `json:"stoploss_current_dist_pct"`
	StoplossCurrentRatio float64      `json:"stoploss_current_dist_ratio"`
	StoplossEntryDist    float64      `json:"stoploss_entry_dist"`
	StoplossEntryRatio   float64      `json:"stoploss_entry_dist_ratio"`
	CurrentRate          float64      `json:"current_rate"`
	TotalProfitAbs       float64      `json:"total_profit_abs"`
	TotalProfitFiat      float64      `json:"total_profit_fiat"`
	TotalProfitRatio     float64      `json:"total_profit_ratio"`
}

type TradeOrder struct {
	Pair                 string   `json:"pair"`
	OrderId              string   `json:"order_id"`
	Status               string   `json:"status"`
	Remaining            float64  `json:"remaining"`
	Amount               float64  `json:"amount"`
	SafePrice            float64  `json:"safe_price"`
	Cost                 float64  `json:"cost"`
	Filled               float64  `json:"filled"`
	FtOrderSide          string   `json:"ft_order_side"`
	OrderType            string   `json:"order_type"`
	IsOpen               bool     `json:"is_open"`
	OrderTimestamp       int64    `json:"order_timestamp"`
	OrderFilledTimestamp int64    `json:"order_filled_timestamp"`
	FtFeeBase            *float64 `json:"ft_fee_base"`
	FtOrderTag           string   `json:"ft_order_tag"`
}

// LoginResponse 登录响应结构
type LoginResponse struct {
	AccessToken  string `json:"access_token"` // JWT 令牌
	RefreshToken string `json:"refresh_token"`
}

type PositionStatus struct {
	Current    int     `json:"current"`     // 当前持仓数量
	Max        int     `json:"max"`         // 最大持仓限制
	TotalStake float64 `json:"total_stake"` // 当前总投入金额
}

type ForceBuyPayload struct {
	Pair      string  `json:"pair"`      // 如 ETH/USDT:USDT
	Price     float64 `json:"price"`     // 限价
	OrderType string  `json:"ordertype"` // "limit" 或 "market"
	Side      string  `json:"side"`      // "long" 或 "short"
	EntryTag  string  `json:"entry_tag"` // 自定义标签，例如 "force_entry"
}

type ForceAdjustBuyPayload struct {
	Pair        string  `json:"pair"`      // 如 ETH/USDT:USDT
	Price       float64 `json:"price"`     // 限价
	OrderType   string  `json:"ordertype"` // "limit" 或 "market"
	Side        string  `json:"side"`      // "long" 或 "short"
	EntryTag    string  `json:"entry_tag"` // 自定义标签，例如 "force_entry"
	StakeAmount float64 `json:"stakeamount"`
}

type ForceSellPayload struct {
	TradeId   string `json:"tradeid"`   // 交易ID
	OrderType string `json:"ordertype"` // "limit" 或 "market"
	Amount    string `json:"amount"`    // 卖出数量
}

// WhitelistResponse whitelist接口响应结构
type WhitelistResponse struct {
	Whitelist []string `json:"whitelist"` // 交易对白名单列表
	Length    int      `json:"length"`    // 白名单长度
	Method    []string `json:"method"`    // 使用的过滤方法
}

// WebhookMessage Freqtrade webhook消息结构
type WebhookMessage struct {
	Type string `json:"type"` // 消息类型：entry, entry_cancel, entry_fill, exit, exit_fill, exit_cancel, status
	// 通用字段
	TradeId       *int     `json:"trade_id,omitempty"`
	Exchange      *string  `json:"exchange,omitempty"`
	Pair          *string  `json:"pair,omitempty"`
	Direction     *string  `json:"direction,omitempty"` // long/short
	Leverage      *float64 `json:"leverage,omitempty"`
	Amount        *float64 `json:"amount,omitempty"`
	StakeAmount   *float64 `json:"stake_amount,omitempty"`
	StakeCurrency *string  `json:"stake_currency,omitempty"`
	BaseCurrency  *string  `json:"base_currency,omitempty"`
	QuoteCurrency *string  `json:"quote_currency,omitempty"`
	FiatCurrency  *string  `json:"fiat_currency,omitempty"`
	OrderType     *string  `json:"order_type,omitempty"`
	CurrentRate   *float64 `json:"current_rate,omitempty"`
	EnterTag      *string  `json:"enter_tag,omitempty"`

	// 价格相关
	OpenRate  *float64 `json:"open_rate,omitempty"`
	CloseRate *float64 `json:"close_rate,omitempty"`
	Limit     *float64 `json:"limit,omitempty"`

	// 时间相关
	OpenDate  *string `json:"open_date,omitempty"`
	CloseDate *string `json:"close_date,omitempty"`

	// 盈亏相关
	Gain         *string  `json:"gain,omitempty"`
	ProfitAmount *float64 `json:"profit_amount,omitempty"`
	ProfitRatio  *float64 `json:"profit_ratio,omitempty"`
	ExitReason   *string  `json:"exit_reason,omitempty"`

	// 状态消息
	Status *string `json:"status,omitempty"`

	// 原始数据（用于灵活处理）
	RawData map[string]interface{} `json:"-"`
}
