package main

import (
	"log"
	"monitor-trade/config"
	"monitor-trade/controller"
	"monitor-trade/controller/binance"
	"monitor-trade/controller/freqtrade"
	"monitor-trade/controller/http"
	"monitor-trade/controller/redis"
	"monitor-trade/controller/tg"
)

func main() {
	conf := config.LoadFromEnv()
	redisController := redis.NewRedisController(conf)

	// 启动时从Redis加载监控数据到本地
	if err := redisController.LoadMonitorPairsFromRedis(); err != nil {
		log.Printf("加载Redis监控数据失败: %v", err)
	}

	// 确保Redis keyspace事件已启用
	if err := redisController.EnableRedisKeyspaceNotifications(); err != nil {
		log.Printf("启用Redis keyspace事件失败: %v", err)
	}

	// 启动Redis keyspace事件监听，自动同步本地数据
	go redisController.StartRedisSync()

	// 初始化Binance控制器
	binanceController := binance.NewBinanceController()

	// 设置 BinanceController 的 RedisController 引用
	binanceController.SetRedisController(redisController)

	// Tg 消息通知通道
	messageChan := make(chan string, 1000)
	freqtradeController := freqtrade.NewFreqtradeController(conf.BotBaseUrl, conf.BotUsername, conf.BotPasswd, redisController)
	freqtradeController.Init(messageChan)

	// 使用Binance作为价格数据源的TgController
	tgController := tg.NewTgController(conf.TelegramToken, conf.TelegramId, redisController, freqtradeController)
	go tgController.SendMessageByChan(messageChan)
	go tgController.HandleCommand()

	mainController := controller.NewMainController(tgController, redisController, conf, freqtradeController, binanceController)
	// 使用Binance WebSocket监听价格变化
	go binanceController.Watch(mainController.WatchKey)
	go mainController.Start()

	httpHandler := http.NewHttpHandler(mainController, redisController, freqtradeController)
	http.ListenAndServe(httpHandler)
}
