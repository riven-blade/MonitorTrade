package model

// BookTicker推送数据结构
type BookTickerData struct {
	EventType       string `json:"e"` // 事件类型 "bookTicker"
	UpdateID        int64  `json:"u"` // 更新ID
	EventTime       int64  `json:"E"` // 事件推送时间
	TransactionTime int64  `json:"T"` // 撮合时间
	Symbol          string `json:"s"` // 交易对
	BidPrice        string `json:"b"` // 买单最优挂单价格
	BidQty          string `json:"B"` // 买单最优挂单数量
	AskPrice        string `json:"a"` // 卖单最优挂单价格
	AskQty          string `json:"A"` // 卖单最优挂单数量
}

// PremiumIndexData 标记价格和资金费率数据结构
type PremiumIndexData struct {
	Symbol               string `json:"symbol"`               // 交易对
	MarkPrice            string `json:"markPrice"`            // 标记价格
	IndexPrice           string `json:"indexPrice"`           // 指数价格
	EstimatedSettlePrice string `json:"estimatedSettlePrice"` // 预估结算价，仅交割合约返回
	LastFundingRate      string `json:"lastFundingRate"`      // 最近更新的资金费率
	NextFundingTime      int64  `json:"nextFundingTime"`      // 下次资金费时间
	InterestRate         string `json:"interestRate"`         // 标的资产基础利率
	Time                 int64  `json:"time"`                 // 更新时间
}

// BinanceAccountInfo 账户信息
type BinanceAccountInfo struct {
	FeeTier                     int        `json:"feeTier"`                     // 手续费等级
	CanTrade                    bool       `json:"canTrade"`                    // 是否可以交易
	CanDeposit                  bool       `json:"canDeposit"`                  // 是否可以入金
	CanWithdraw                 bool       `json:"canWithdraw"`                 // 是否可以出金
	UpdateTime                  int64      `json:"updateTime"`                  // 保留字段，请忽略
	TotalInitialMargin          string     `json:"totalInitialMargin"`          // 当前所需起始保证金总额(存在逐仓请忽略), 仅计算usdt资产
	TotalMaintMargin            string     `json:"totalMaintMargin"`            // 维持保证金总额, 仅计算usdt资产
	TotalWalletBalance          string     `json:"totalWalletBalance"`          // 账户余额总额, 仅计算usdt资产
	TotalUnrealizedProfit       string     `json:"totalUnrealizedProfit"`       // 持仓未实现盈亏总额, 仅计算usdt资产
	TotalMarginBalance          string     `json:"totalMarginBalance"`          // 保证金总额, 仅计算usdt资产
	TotalPositionInitialMargin  string     `json:"totalPositionInitialMargin"`  // 持仓所需起始保证金(基于最新标记价格), 仅计算usdt资产
	TotalOpenOrderInitialMargin string     `json:"totalOpenOrderInitialMargin"` // 当前挂单所需起始保证金(基于最新标记价格), 仅计算usdt资产
	TotalCrossWalletBalance     string     `json:"totalCrossWalletBalance"`     // 全仓账户余额, 仅计算usdt资产
	TotalCrossUnPnl             string     `json:"totalCrossUnPnl"`             // 全仓持仓未实现盈亏总额, 仅计算usdt资产
	AvailableBalance            string     `json:"availableBalance"`            // 可用余额, 仅计算usdt资产
	MaxWithdrawAmount           string     `json:"maxWithdrawAmount"`           // 最大可转出余额, 仅计算usdt资产
	Assets                      []Asset    `json:"assets"`                      // 资产内容
	Positions                   []Position `json:"positions"`                   // 持仓信息
}

// Asset 资产信息
type Asset struct {
	Asset                  string `json:"asset"`                  // 资产
	WalletBalance          string `json:"walletBalance"`          // 余额
	UnrealizedProfit       string `json:"unrealizedProfit"`       // 未实现盈亏
	MarginBalance          string `json:"marginBalance"`          // 保证金余额
	MaintMargin            string `json:"maintMargin"`            // 维持保证金
	InitialMargin          string `json:"initialMargin"`          // 当前所需起始保证金
	PositionInitialMargin  string `json:"positionInitialMargin"`  // 持仓所需起始保证金(基于最新标记价格)
	OpenOrderInitialMargin string `json:"openOrderInitialMargin"` // 当前挂单所需起始保证金(基于最新标记价格)
	CrossWalletBalance     string `json:"crossWalletBalance"`     // 全仓账户余额
	CrossUnPnl             string `json:"crossUnPnl"`             // 全仓持仓未实现盈亏
	AvailableBalance       string `json:"availableBalance"`       // 可用余额
	MaxWithdrawAmount      string `json:"maxWithdrawAmount"`      // 最大可转出余额
	MarginAvailable        bool   `json:"marginAvailable"`        // 是否可用作联合保证金
	UpdateTime             int64  `json:"updateTime"`             // 更新时间
}

// Position 持仓信息
type Position struct {
	Symbol                 string `json:"symbol"`                 // 交易对
	InitialMargin          string `json:"initialMargin"`          // 当前所需起始保证金(基于最新标记价格)
	MaintMargin            string `json:"maintMargin"`            // 维持保证金
	UnrealizedProfit       string `json:"unrealizedProfit"`       // 持仓未实现盈亏
	PositionInitialMargin  string `json:"positionInitialMargin"`  // 持仓所需起始保证金(基于最新标记价格)
	OpenOrderInitialMargin string `json:"openOrderInitialMargin"` // 当前挂单所需起始保证金(基于最新标记价格)
	Leverage               string `json:"leverage"`               // 杠杆倍率
	Isolated               bool   `json:"isolated"`               // 是否是逐仓模式
	EntryPrice             string `json:"entryPrice"`             // 持仓成本价
	MarkPrice              string `json:"markPrice"`              // 当前标记价格
	MaxNotional            string `json:"maxNotional"`            // 当前杠杆下用户可用的最大名义价值
	PositionSide           string `json:"positionSide"`           // 持仓方向
	PositionAmt            string `json:"positionAmt"`            // 持仓数量
	Notional               string `json:"notional"`               // 持仓名义价值
	IsolatedWallet         string `json:"isolatedWallet"`         // 逐仓保证金
	UpdateTime             int64  `json:"updateTime"`             // 更新时间
	BidNotional            string `json:"bidNotional"`            // 买单净值
	AskNotional            string `json:"askNotional"`            // 卖单净值
}
