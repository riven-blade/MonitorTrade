package binance

import (
	"context"
	"fmt"
	"log"
	"monitor-trade/controller/redis"
	"monitor-trade/model"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

type BinanceController struct {
	conn               *websocket.Conn
	redisController    *redis.RedisController
	changePairDataChan chan model.PairData
	ctx                context.Context
	cancel             context.CancelFunc
	httpClient         *http.Client // HTTP客户端用于REST API调用
}

func NewBinanceController() *BinanceController {
	ctx, cancel := context.WithCancel(context.Background())
	return &BinanceController{
		changePairDataChan: make(chan model.PairData, 1000),
		ctx:                ctx,
		cancel:             cancel,
		httpClient:         &http.Client{Timeout: 10 * time.Second},
	}
}

// SetRedisController 设置 Redis 控制器
func (b *BinanceController) SetRedisController(redisController *redis.RedisController) {
	b.redisController = redisController
}

// 连接到Binance WebSocket推送流
func (b *BinanceController) Connect() error {
	// 使用期货合约的全市场最优挂单信息流
	wsURL := "wss://fstream.binance.com/ws/!bookTicker"

	u, err := url.Parse(wsURL)
	if err != nil {
		return fmt.Errorf("解析WebSocket URL失败: %v", err)
	}

	log.Printf("正在连接到Binance期货合约最优挂单推送流 (每5秒推送): %s", wsURL)

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return fmt.Errorf("连接Binance期货WebSocket推送流失败: %v", err)
	}

	b.conn = conn
	log.Println("成功连接到Binance期货合约最优挂单推送流")

	return nil
}

// 启动价格监听
func (b *BinanceController) Watch(changePairDataChan chan model.PairData) {
	b.changePairDataChan = changePairDataChan

	if err := b.Connect(); err != nil {
		log.Printf("连接Binance失败: %v", err)
		return
	}

	defer b.conn.Close()

	// 启动重连机制
	go b.handleReconnect()

	log.Println("开始监听Binance最优挂单推送流...")

	for {
		select {
		case <-b.ctx.Done():
			return
		default:
			var ticker model.BookTickerData
			err := b.conn.ReadJSON(&ticker)
			if err != nil {
				log.Printf("读取Binance推送数据失败: %v", err)
				time.Sleep(time.Second)
				continue
			}

			go b.processBookTicker(ticker)
		}
	}
}

// 处理BookTicker推送数据
func (b *BinanceController) processBookTicker(ticker model.BookTickerData) {
	// ticker pair 以USDT结尾的交易对
	if !strings.HasSuffix(ticker.Symbol, "USDT") {
		return
	}

	// 转换符号格式，从BTCUSDT到BTC/USDT
	pair := b.formatPairSymbol(ticker.Symbol)

	// 解析价格数据
	pairData, err := b.convertBookTickerToPairData(ticker)
	if err != nil {
		log.Printf("转换价格数据失败 %s: %v", pair, err)
		return
	}

	// 更新价格数据（RedisController 内部有锁保护）- 所有数据都存储
	if b.redisController != nil {
		b.redisController.UpdatePairPrice(pair, pairData)
	}

	// 如果设置了监听列表，只向 WatchKey channel 发送列表中的交易对
	if b.redisController != nil && !b.redisController.IsWatchedPair(pair) {
		return // 跳过不在监听列表中的交易对，但数据已经存储了
	}

	// 通知价格更新
	go func() {
		select {
		case b.changePairDataChan <- *pairData:
		default:
			// 如果channel满了，跳过这次通知
			log.Printf("价格更新通知channel满了，跳过这次通知: %s", pair)
		}
	}()
}

// 转换Binance符号格式
func (b *BinanceController) formatPairSymbol(symbol string) string {
	// 将BTCUSDT转换为BTC/USDT:USDT格式
	if len(symbol) >= 4 && symbol[len(symbol)-4:] == "USDT" {
		base := symbol[:len(symbol)-4]
		return fmt.Sprintf("%s/USDT:USDT", base)
	}
	if len(symbol) >= 3 && symbol[len(symbol)-3:] == "BTC" {
		base := symbol[:len(symbol)-3]
		return fmt.Sprintf("%s/BTC:BTC", base)
	}
	if len(symbol) >= 3 && symbol[len(symbol)-3:] == "ETH" {
		base := symbol[:len(symbol)-3]
		return fmt.Sprintf("%s/ETH:ETH", base)
	}
	if len(symbol) >= 3 && symbol[len(symbol)-4:] == "USDC" {
		base := symbol[:len(symbol)-4]
		return fmt.Sprintf("%s/USDC:USDC", base)
	}
	return symbol
}

// 转换BookTicker数据为PairData格式
func (b *BinanceController) convertBookTickerToPairData(ticker model.BookTickerData) (*model.PairData, error) {
	var pairData model.PairData

	// 解析买单最优挂单价格
	if bidPrice, err := strconv.ParseFloat(ticker.BidPrice, 64); err == nil {
		pairData.BidPrice = bidPrice
	} else {
		return nil, fmt.Errorf("无法解析买单价格: %s", ticker.BidPrice)
	}

	// 解析卖单最优挂单价格
	if askPrice, err := strconv.ParseFloat(ticker.AskPrice, 64); err == nil {
		pairData.AskPrice = askPrice
	} else {
		return nil, fmt.Errorf("无法解析卖单价格: %s", ticker.AskPrice)
	}

	pairData.Close = (pairData.AskPrice + pairData.BidPrice) / 2
	pairData.Pair = b.formatPairSymbol(ticker.Symbol)
	// 使用 ticker 中的撮合时间戳，将毫秒转换为秒
	if ticker.TransactionTime > 0 {
		pairData.Timestamp = time.Unix(ticker.TransactionTime/1000, 0).Format("2006-01-02 15:04:05")
	} else {
		// 如果时间戳无效，回退到当前时间
		pairData.Timestamp = time.Now().Format("2006-01-02 15:04:05")
	}
	return &pairData, nil
}

// 处理重连
func (b *BinanceController) handleReconnect() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-b.ctx.Done():
			return
		case <-ticker.C:
			// 发送ping保持连接
			if b.conn != nil {
				if err := b.conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(time.Second)); err != nil {
					log.Printf("发送ping失败: %v", err)
					b.reconnect()
				}
			}
		}
	}
}

// 重新连接
func (b *BinanceController) reconnect() {
	log.Println("尝试重新连接Binance WebSocket...")

	if b.conn != nil {
		b.conn.Close()
	}

	for i := 0; i < 5; i++ {
		if err := b.Connect(); err != nil {
			log.Printf("重连失败 (attempt %d/5): %v", i+1, err)
			time.Sleep(time.Duration(i+1) * time.Second)
			continue
		}
		log.Println("重连成功")
		return
	}

	log.Println("重连失败，将在60秒后重试")
	time.Sleep(60 * time.Second)
	go b.reconnect()
}

// 停止监听
func (b *BinanceController) Stop() {
	if b.cancel != nil {
		b.cancel()
	}
	if b.conn != nil {
		b.conn.Close()
	}
}
