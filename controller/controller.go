package controller

import (
	"fmt"
	"log"
	"monitor-trade/config"
	"monitor-trade/controller/binance"
	"monitor-trade/controller/redis"
	"monitor-trade/controller/tg"
	"monitor-trade/model"
)

type MainController struct {
	TgController      *tg.TgController
	RedisController   *redis.RedisController
	Conf              *config.Config
	BinanceController *binance.BinanceController
	WatchKey          chan model.PairData
	TradeChan         chan model.ForceBuyPayload
}

// NewMainController 创建MainController
func NewMainController(tgController *tg.TgController, redisController *redis.RedisController,
	conf *config.Config, binanceController *binance.BinanceController, tradeChan chan model.ForceBuyPayload) *MainController {
	return &MainController{
		TgController:      tgController,
		RedisController:   redisController,
		WatchKey:          make(chan model.PairData, 200),
		Conf:              conf,
		BinanceController: binanceController,
		TradeChan:         tradeChan,
	}
}

func (c *MainController) Start() {
	for pairData := range c.WatchKey {
		// 处理短线
		go c.HandleShort(&pairData)
		// 处理长线
		go c.HandleLong(&pairData)
	}
}

func (c *MainController) HandleShort(pairData *model.PairData) {
	shortData, exists := c.RedisController.GetMonitorPair(pairData.Pair, tg.ShortDirect)
	if !exists {
		return
	}
	if shortData.Price <= 0 {
		return
	}

	if pairData.AskPrice > shortData.Price {
		// 检查资金费率条件
		fundingRate, err := c.BinanceController.GetFundingRate(pairData.Pair)
		if err != nil {
			log.Printf("获取交易对 %s 的资金费率失败: %v", pairData.Pair, err)
			resultMsg := fmt.Sprintf("❌ %s 做空操作失败: 获取资金费率失败 - %v", pairData.Pair, err)
			c.TgController.SendMessage(resultMsg)
			return
		}
		if fundingRate <= c.Conf.FundingRate {
			log.Printf("%s 的资金费率 %.6f 小于等于阈值 %.6f，跳过做空处理",
				pairData.Pair, fundingRate, c.Conf.FundingRate)
			return
		}

		log.Printf("时间戳 %s 交易对 %s 的当前卖单价 %.6f 高于做空价格 %.6f，资金费率 %.6f，执行做空操作",
			pairData.Timestamp, pairData.Pair, pairData.AskPrice, shortData.Price, fundingRate)

		c.TradeChan <- model.ForceBuyPayload{
			Pair:      pairData.Pair,
			Price:     pairData.AskPrice,
			Side:      "short",
			EntryTag:  "force_buy",
			OrderType: "limit",
		}
	}
}

func (c *MainController) HandleLong(pairData *model.PairData) {
	longData, exists := c.RedisController.GetMonitorPair(pairData.Pair, tg.LongDirect)
	if !exists {
		return
	}
	if longData.Price <= 0 {
		return
	}
	if pairData.BidPrice < longData.Price {
		// 当前买单价(最低价) < 做多价格
		log.Printf("时间戳 %s 交易对 %s 的当前买单价 %.6f 低于做多价格 %.6f，执行做多操作",
			pairData.Timestamp, pairData.Pair, pairData.BidPrice, longData.Price)

		c.TradeChan <- model.ForceBuyPayload{
			Pair:      pairData.Pair,
			Price:     pairData.BidPrice,
			Side:      "long",
			EntryTag:  "force_buy",
			OrderType: "limit",
		}
	}
}
