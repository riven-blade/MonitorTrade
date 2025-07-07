package tg

import (
	"fmt"
	"log"
	"monitor-trade/model"
	"strings"
)

// 处理 /short 命令
func (tg *TgController) handleShortCommand(pair string, price float64) string {
	data, _ := tg.RedisController.GetMonitorPair(pair, ShortDirect)
	data.Pair = pair
	resultMsg := ""

	dataPair := tg.RedisController.GetPairPrice(data.Pair)
	// 计算中间价作为当前价格
	currentPrice := (dataPair.BidPrice + dataPair.AskPrice) / 2
	if currentPrice <= 0 {
		return fmt.Sprintf("❌ 无法获取 %s 的最新价格，请检查交易对是否存在", pair)
	}
	if currentPrice > price {
		return fmt.Sprintf("❌ 当前价格 %.6f 大于设置的限价 %.6f，请调整限价", currentPrice, price)
	}

	if data.Price > 0 {
		oldPrice := data.Price
		data.Price = price
		resultMsg = fmt.Sprintf("🟢 %s 做空监听，新限价: %.6f，旧限价: %.6f", pair, data.Price, oldPrice)
	} else {
		data.Price = price
		resultMsg = fmt.Sprintf("🟢 %s 做空监听，限价: %.6f", pair, data.Price)
	}

	if err := tg.RedisController.SetMonitorPair(data, ShortDirect); err != nil {
		resultMsg = fmt.Sprintf("设置 %s 做空监听失败: %v", pair, err)
	} else {
		resultMsg += fmt.Sprintf(", 当前价格: %.6f", currentPrice)
	}
	return resultMsg
}

// 处理 /long 命令
func (tg *TgController) handleLongCommand(pair string, price float64) string {
	data, _ := tg.RedisController.GetMonitorPair(pair, LongDirect)
	data.Pair = pair

	// 获取当前交易对的最新价格
	dataPair := tg.RedisController.GetPairPrice(data.Pair)
	// 计算中间价作为当前价格
	currentPrice := (dataPair.BidPrice + dataPair.AskPrice) / 2
	if currentPrice <= 0 {
		return fmt.Sprintf("❌ 无法获取 %s 的最新价格，请检查交易对是否存在", pair)
	}
	if currentPrice < price {
		return fmt.Sprintf("❌ 当前价格 %.6f 小于设置的限价 %.6f，请调整限价", currentPrice, price)
	}

	resultMsg := ""
	if data.Price > 0 {
		oldPrice := data.Price
		data.Price = price
		resultMsg = fmt.Sprintf("🟢 %s 做多监听，新限价: %.6f，旧限价: %.6f", pair, data.Price, oldPrice)
	} else {
		data.Price = price
		resultMsg = fmt.Sprintf("🟢 %s 做多监听，限价: %.6f", pair, data.Price)
	}
	if err := tg.RedisController.SetPairDataToRedis(data, LongDirect); err != nil {
		resultMsg = fmt.Sprintf("设置 %s 做多监听失败: %v", pair, err)
	} else {
		resultMsg += fmt.Sprintf(", 当前价格: %.6f", currentPrice)
	}
	return resultMsg
}

// 处理 /cancel 命令
func (tg *TgController) handleCancelCommand(pair string, direct string) string {
	resultMsg := ""
	switch direct {
	case LongDirect:
		resultMsg = fmt.Sprintf("✅ %s 做多监听已取消", pair)
		log.Printf("取消 %s 做多监听", pair)
	case ShortDirect:
		resultMsg = fmt.Sprintf("✅ %s 做空监听已取消", pair)
		log.Printf("取消 %s 做空监听", pair)
	default:
		resultMsg = fmt.Sprintf("❌ 无效的方向: %s。请使用 'long' 或 'short'", direct)
		log.Printf("无效的方向: %s", direct)
		return resultMsg
	}
	tg.RedisController.DeleteMonitorPair(pair, direct)
	return resultMsg
}

func (tg *TgController) handleShowConfigCommand() string {
	resultMsg := ""
	// 查找所有交易对
	pairsData, err := tg.RedisController.GetAllPairPricesData()
	if err != nil {
		resultMsg = fmt.Sprintf("获取交易对数据失败: %v", err)
		return resultMsg
	}

	// 查看被监听的交易对
	monitorShortPairsData := tg.RedisController.GetAllMonitorPairsData("short")
	monitorShortPair := make(map[string]float64, 300)
	unMonitoredShortPairs := make(map[string]float64, 300)

	for i := range monitorShortPairsData {
		monitorShortPair[monitorShortPairsData[i].PairMonitorData.Pair] = float64(monitorShortPairsData[i].TTL)
	}
	for i := range pairsData {
		pair := pairsData[i].Pair
		if _, ok := monitorShortPair[pair]; !ok {
			unMonitoredShortPairs[pair] = 0
		}
	}

	resultMsg += "Short:\n"
	resultMsg += "Monitor Pair:\n"
	for i := range monitorShortPair {
		if strings.HasSuffix(i, "/USDT:USDT") {
			resultMsg += fmt.Sprintf("%s ", strings.TrimSuffix(i, "/USDT:USDT"))
		} else {
			resultMsg += fmt.Sprintf("%s ", i)
		}
	}

	resultMsg += "\n\n"
	resultMsg += "Expiration Pair:\n"
	for i := range monitorShortPair {
		// 只显示过期时间小于一天的交易对
		if monitorShortPair[i] < 86400 {
			if strings.HasSuffix(i, "/USDT:USDT") {
				resultMsg += fmt.Sprintf("%s ", strings.TrimSuffix(i, "/USDT:USDT"))
			} else {
				resultMsg += fmt.Sprintf("%s ", i)
			}
		}
	}

	resultMsg += "\n\n"
	resultMsg += "Unmonitored Pair:\n"
	for i := range unMonitoredShortPairs {
		if strings.HasSuffix(i, "/USDT:USDT") {
			resultMsg += fmt.Sprintf("%s ", strings.TrimSuffix(i, "/USDT:USDT"))
		} else {
			resultMsg += fmt.Sprintf("%s ", i)
		}
	}

	// 查看被监听的交易对
	monitorLongPairsData := tg.RedisController.GetAllMonitorPairsData("long")
	monitorLongPair := make(map[string]float64, 300)
	unMonitoredLongPairs := make(map[string]float64, 300)

	for i := range monitorLongPairsData {
		monitorLongPair[monitorLongPairsData[i].PairMonitorData.Pair] = float64(monitorLongPairsData[i].TTL)
	}
	for i := range pairsData {
		pair := pairsData[i].Pair
		if _, ok := monitorLongPair[pair]; !ok {
			unMonitoredLongPairs[pair] = 0
		}
	}

	resultMsg += "\n\n"
	resultMsg += "Long:\n"
	resultMsg += "Monitor Pair:\n"
	for i := range monitorLongPair {
		if strings.HasSuffix(i, "/USDT:USDT") {
			resultMsg += fmt.Sprintf("%s ", strings.TrimSuffix(i, "/USDT:USDT"))
		} else {
			resultMsg += fmt.Sprintf("%s ", i)
		}
	}

	resultMsg += "\n\n"
	resultMsg += "Expiration Pair:\n"
	for i := range monitorLongPair {
		// 只显示过期时间小于一天的交易对
		if monitorLongPair[i] < 86400 {
			if strings.HasSuffix(i, "/USDT:USDT") {
				resultMsg += fmt.Sprintf("%s ", strings.TrimSuffix(i, "/USDT:USDT"))
			} else {
				resultMsg += fmt.Sprintf("%s ", i)
			}
		}
	}
	return resultMsg
}

func (tg *TgController) handleShowCommand(pair string) string {
	resultMsg := ""
	// 查找所有交易对
	pairsData := tg.RedisController.GetPairPrice(pair)

	// 查看被监听的交易对
	monitorLongData, _ := tg.RedisController.GetMonitorPair(pair, LongDirect)
	monitorShortData, _ := tg.RedisController.GetMonitorPair(pair, ShortDirect)

	if monitorLongData.Price > 0 {
		resultMsg += fmt.Sprintf("%s 做多监听，限价: %.6f\n", pair, monitorLongData.Price)
	}
	if monitorShortData.Price > 0 {
		resultMsg += fmt.Sprintf("%s 做空监听，限价: %.6f\n", pair, monitorShortData.Price)
	}
	// 计算中间价作为当前价格
	currentPrice := (pairsData.BidPrice + pairsData.AskPrice) / 2
	resultMsg += fmt.Sprintf("当前价格: %.6f\n", currentPrice)
	return resultMsg
}

func (tg *TgController) handleWhiteList() string {
	resultMsg := ""
	// 查找所有交易对
	whiteListPairs := tg.RedisController.WatchedPairs
	resultMsg += ""
	for i := range whiteListPairs {
		pairData := whiteListPairs[i]
		if strings.HasSuffix(pairData, "/USDT:USDT") {
			resultMsg += fmt.Sprintf("%s ", strings.TrimSuffix(pairData, "/USDT:USDT"))
		} else {
			resultMsg += fmt.Sprintf("%s ", pairData)
		}
	}
	return resultMsg
}

// 处理 /ad 命令
func (tg *TgController) handleADCommand(pair string, stakeAmount float64, price float64) string {
	// 获取当前交易对的最新价格
	dataPair := tg.RedisController.GetPairPrice(pair)
	if dataPair.BidPrice <= 0 || dataPair.AskPrice <= 0 {
		return fmt.Sprintf("❌ 无法获取 %s 的最新价格，请检查交易对是否存在", pair)
	}

	hasTrade := false
	isShort := false
	tradeStatus := tg.FreqtradeController.TradeStatus
	for i := range tradeStatus {
		if tradeStatus[i].Pair == pair {
			hasTrade = true
			isShort = tradeStatus[i].IsShort
		}
	}

	if !hasTrade {
		return fmt.Sprintf("❌ %s 没有持仓", pair)
	}

	if stakeAmount <= 0 {
		stakeAmount = 10
	}

	// 如果没有指定价格，使用当前市价
	if price <= 0 {
		if isShort {
			price = dataPair.AskPrice
		} else {
			price = dataPair.BidPrice
		}
	}

	if isShort {
		err := tg.FreqtradeController.ForceAdjustBuy(pair, price, ShortDirect, stakeAmount)
		if err != nil {
			log.Printf("%s 做空加仓失败: %v", pair, err)
			return fmt.Sprintf("❌ %s 做空加仓失败: %v", pair, err)
		}
		return fmt.Sprintf("📉 %s 做空加仓成功，金额: %.2f，价格: %.6f", pair, stakeAmount, price)
	} else {
		err := tg.FreqtradeController.ForceAdjustBuy(pair, price, LongDirect, stakeAmount)
		if err != nil {
			log.Printf("%s 做多加仓失败: %v", pair, err)
			return fmt.Sprintf("❌ %s 做多加仓失败: %v", pair, err)
		}
		return fmt.Sprintf("📈 %s 做多加仓成功，金额: %.2f，价格: %.6f", pair, stakeAmount, price)
	}
}

// 处理 /pc 命令（平仓）
func (tg *TgController) handlePCCommand(pair string, amount float64) string {
	// 获取当前交易状态
	tradeStatus := tg.FreqtradeController.TradeStatus

	var targetTrade *model.TradePosition
	for i := range tradeStatus {
		if tradeStatus[i].Pair == pair {
			targetTrade = &tradeStatus[i]
			break
		}
	}

	if targetTrade == nil {
		return fmt.Sprintf("❌ %s 没有开仓", pair)
	}

	// 参数验证
	if amount <= 0 {
		return "❌ 平仓金额必须大于0"
	}

	if targetTrade.StakeAmount <= 0 {
		return fmt.Sprintf("❌ %s 投入金额数据异常", pair)
	}

	// 计算平仓的amount
	pcStakeAmount := amount // 用户想要平仓的投入金额
	pcRate := pcStakeAmount / targetTrade.StakeAmount
	if pcRate > 0.9 {
		pcRate = 0.9
	}
	amount = targetTrade.Amount * pcRate

	// 使用 ForceSell 进行平仓
	tradeIdStr := fmt.Sprintf("%d", targetTrade.TradeId)
	amountStr := fmt.Sprintf("%.6f", amount)

	err := tg.FreqtradeController.ForceSell(tradeIdStr, "market", amountStr)
	if err != nil {
		log.Printf("%s 平仓失败: %v", pair, err)
		return fmt.Sprintf("❌ %s 平仓失败: %v", pair, err)
	}

	if targetTrade.IsShort {
		return fmt.Sprintf("📈 %s 做空平仓成功，数量: %.6f, stake: %.2f", pair, amount, pcStakeAmount)
	} else {
		return fmt.Sprintf("📉 %s 做多平仓成功，数量: %.6f, stake: %.2f", pair, amount, pcStakeAmount)
	}
}

// 处理 /adjust 命令（无参数时显示仓位信息）
func (tg *TgController) handleShowPositionsCommand() string {
	resultMsg := ""

	// 获取 Freqtrade 实际交易状态
	tradeStatus := tg.FreqtradeController.TradeStatus

	if len(tradeStatus) == 0 {
		resultMsg += "无仓位\n"
		return resultMsg
	}

	// 按做多/做空分类显示
	longPositions := []model.TradePosition{}
	shortPositions := []model.TradePosition{}

	for i := range tradeStatus {
		trade := tradeStatus[i]
		if trade.IsOpen {
			if trade.IsShort {
				shortPositions = append(shortPositions, trade)
			} else {
				longPositions = append(longPositions, trade)
			}
		}
	}

	// 显示做多仓位
	if len(longPositions) > 0 {
		resultMsg += "📈 **做多仓位:**\n"
		for i := range longPositions {
			trade := longPositions[i]
			if len(trade.Orders) > 0 {
				resultMsg += fmt.Sprintf("%s %.2f\n", trade.Pair, trade.Orders[0].Cost/trade.Leverage)
			} else {
				resultMsg += fmt.Sprintf("%s\n", trade.Pair)
			}
		}
		resultMsg += "\n"
	}

	// 显示做空仓位
	if len(shortPositions) > 0 {
		resultMsg += "📉 **做空仓位:**\n"
		for i := range shortPositions {
			trade := shortPositions[i]
			if len(trade.Orders) > 0 {
				resultMsg += fmt.Sprintf("%s %.2f\n", trade.Pair, trade.Orders[0].Cost/trade.Leverage)
			} else {
				resultMsg += fmt.Sprintf("%s\n", trade.Pair)
			}
		}
		resultMsg += "\n"
	}
	return resultMsg
}
