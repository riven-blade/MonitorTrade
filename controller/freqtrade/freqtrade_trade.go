package freqtrade

import (
	"context"
	"fmt"
	"log"
	"monitor-trade/model"
)

// HandleTradeChan 处理交易通道，支持优雅停止
func (fc *FreqtradeController) HandleTradeChan(ctx context.Context, tradeChan chan model.ForceBuyPayload) {
	log.Println("交易处理器已启动")
	defer log.Println("交易处理器已停止")

	for {
		select {
		case <-ctx.Done():
			log.Println("收到停止信号，交易处理器正在停止...")
			return
		case trade := <-tradeChan:
			fc.processTrade(trade)
		}
	}
}

// processTrade 处理单个交易请求
func (fc *FreqtradeController) processTrade(trade model.ForceBuyPayload) {
	log.Printf("收到%s交易请求: %s, 价格: %.6f", trade.Side, trade.Pair, trade.Price)

	// 尝试获取Redis分布式锁
	if !fc.redisController.AcquireTradeLock(trade.Pair) {
		log.Printf("⏰ %s 交易锁获取失败，可能有其他交易正在进行，跳过执行", trade.Pair)
		return
	}

	log.Printf("🔒 获取 %s 交易锁成功，开始处理交易", trade.Pair)

	// 校验仓位限制
	if !fc.CheckForceBuy(trade.Pair) {
		errMsg := fmt.Sprintf("交易对 %s 校验仓位不通过，跳过%s操作", trade.Pair, trade.Side)
		log.Printf("❌ %s", errMsg)
		fc.sendTradeResult(trade.Pair, trade.Price, trade.Side, fmt.Errorf("仓位校验失败"))
		return
	}

	// 执行交易
	err := fc.ForceBuy(trade)
	if err != nil {
		log.Printf("❌ %s %s操作失败: %v", trade.Pair, trade.Side, err)
	} else {
		log.Printf("✅ %s %s操作提交成功，价格: %.6f", trade.Pair, trade.Side, trade.Price)
	}

	// 异步发送结果通知
	go fc.sendTradeResult(trade.Pair, trade.Price, trade.Side, err)
}

// sendTradeResult 统一处理交易结果的消息发送
func (fc *FreqtradeController) sendTradeResult(pair string, price float64, side string, err error) {
	var resultMsg string
	if err != nil {
		resultMsg = fmt.Sprintf("❌ %s %s操作失败: %v", pair, side, err)
		// 交易失败时删除Redis中的监控数据
		fc.redisController.DeleteMonitorPair(pair, side)
	} else {
		resultMsg = fmt.Sprintf("✅ %s %s操作提交成功，价格: %.6f", pair, side, price)
	}

	// 安全地发送消息，避免阻塞
	select {
	case fc.messageChan <- resultMsg:
		// 消息发送成功
	default:
		log.Printf("⚠️ 消息通道已满，跳过发送: %s", resultMsg)
	}
}
