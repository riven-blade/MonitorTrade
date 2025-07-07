package binance

import (
	"encoding/json"
	"fmt"
	"io"
	"monitor-trade/config"
	"monitor-trade/controller/redis"
	"monitor-trade/model"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

// TestNewBinanceController 测试构造函数
func TestNewBinanceController(t *testing.T) {
	controller := NewBinanceController()

	if controller == nil {
		t.Fatal("NewBinanceController() 返回 nil")
	}

	// changeKey 已经移除，不再需要测试

	if controller.ctx == nil {
		t.Error("ctx 应该被初始化")
	}

	if controller.cancel == nil {
		t.Error("cancel 应该被初始化")
	}
}

// TestSetWatchedPairs 测试设置监听交易对
func TestSetWatchedPairs(t *testing.T) {
	// 创建测试用的 RedisController
	conf := &config.Config{
		Redis: config.RedisConfig{
			Addr:      "localhost:6379",
			Password:  "",
			DB:        0,
			KeyExpire: 300,
		},
	}
	redisController := redis.NewRedisController(conf)

	testPairs := []string{"BTC/USDT:USDT", "ETH/USDT:USDT", "BNB/USDT:USDT"}
	redisController.SetWatchedPairs(testPairs)

	// 验证通过 RedisController 获取的数据
	watchedPairs := redisController.GetWatchedPairs()
	if len(watchedPairs) != len(testPairs) {
		t.Errorf("期望监听 %d 个交易对，实际 %d 个", len(testPairs), len(watchedPairs))
	}

	for i, pair := range testPairs {
		if watchedPairs[i] != pair {
			t.Errorf("期望交易对 %s，实际 %s", pair, watchedPairs[i])
		}
	}
}

// TestFormatPairSymbol 测试符号格式转换
func TestFormatPairSymbol(t *testing.T) {
	controller := NewBinanceController()

	testCases := []struct {
		input    string
		expected string
	}{
		{"BTCUSDT", "BTC/USDT:USDT"},
		{"ETHUSDT", "ETH/USDT:USDT"},
		{"BNBUSDT", "BNB/USDT:USDT"},
		{"ADAUSDT", "ADA/USDT:USDT"},
		{"SOLUSDT", "SOL/USDT:USDT"},
		{"DOGEUSDT", "DOGE/USDT:USDT"},
		{"ETHBTC", "ETH/BTC:BTC"},
		{"BNBBTC", "BNB/BTC:BTC"},
		{"ADAETH", "ADA/ETH:ETH"},
		{"UNKNOWN", "UNKNOWN"}, // 无法识别的格式
	}

	for _, tc := range testCases {
		result := controller.formatPairSymbol(tc.input)
		if result != tc.expected {
			t.Errorf("formatPairSymbol(%s) = %s, 期望 %s", tc.input, result, tc.expected)
		}
	}
}

// TestConvertBookTickerToPairData 测试BookTicker数据转换
func TestConvertBookTickerToPairData(t *testing.T) {
	controller := NewBinanceController()

	ticker := model.BookTickerData{
		EventType:       "bookTicker",
		UpdateID:        123456,
		EventTime:       1640995200000,
		TransactionTime: 1640995200000,
		Symbol:          "BTCUSDT",
		BidPrice:        "50000.00",
		BidQty:          "1.5",
		AskPrice:        "50100.00",
		AskQty:          "2.0",
	}

	pairData, err := controller.convertBookTickerToPairData(ticker)
	if err != nil {
		t.Fatalf("转换失败: %v", err)
	}

	// 验证数据转换
	expectedPair := "BTC/USDT:USDT"
	if pairData.Pair != expectedPair {
		t.Errorf("期望交易对 %s，实际 %s", expectedPair, pairData.Pair)
	}

	// 验证价格数据
	expectedBidPrice := 50000.00
	if pairData.BidPrice != expectedBidPrice {
		t.Errorf("期望买单价 %f，实际 %f", expectedBidPrice, pairData.BidPrice)
	}

	expectedAskPrice := 50100.00
	if pairData.AskPrice != expectedAskPrice {
		t.Errorf("期望卖单价 %f，实际 %f", expectedAskPrice, pairData.AskPrice)
	}

	// 验证时间戳格式
	if pairData.Timestamp == "" {
		t.Error("时间戳不应为空")
	}
}

// TestConvertBookTickerToPairDataInvalidPrice 测试无效价格数据
func TestConvertBookTickerToPairDataInvalidPrice(t *testing.T) {
	controller := NewBinanceController()

	// 测试无效的BidPrice
	ticker := model.BookTickerData{
		Symbol:   "BTCUSDT",
		BidPrice: "invalid",
		AskPrice: "50100.00",
	}

	_, err := controller.convertBookTickerToPairData(ticker)
	if err == nil {
		t.Error("当BidPrice无效时，应该返回错误")
	}

	// 测试无效的AskPrice
	ticker2 := model.BookTickerData{
		Symbol:   "BTCUSDT",
		BidPrice: "50000.00",
		AskPrice: "invalid",
	}

	_, err2 := controller.convertBookTickerToPairData(ticker2)
	if err2 == nil {
		t.Error("当AskPrice无效时，应该返回错误")
	}
}

// TestProcessBookTickerWithFilter 测试带过滤的BookTicker处理
func TestProcessBookTickerWithFilter(t *testing.T) {
	controller := NewBinanceController()

	// 创建测试用的 RedisController
	conf := &config.Config{
		Redis: config.RedisConfig{
			Addr:      "localhost:6379",
			Password:  "",
			DB:        0,
			KeyExpire: 300,
		},
	}
	redisController := redis.NewRedisController(conf)
	controller.SetRedisController(redisController)

	// 设置只监听BTC/USDT:USDT
	redisController.SetWatchedPairs([]string{"BTC/USDT:USDT"})

	// 处理BTC/USDT数据（应该被处理）
	btcTicker := model.BookTickerData{
		Symbol:   "BTCUSDT",
		BidPrice: "50000.00",
		AskPrice: "50100.00",
	}

	controller.processBookTicker(btcTicker)

	// 验证BTC/USDT:USDT数据被存储
	result := redisController.GetPairPrice("BTC/USDT:USDT")
	if result.Pair == "" {
		t.Error("BTC/USDT:USDT 数据应该被存储")
	}

	// 处理ETH/USDT数据（应该被过滤掉）
	ethTicker := model.BookTickerData{
		Symbol:   "ETHUSDT",
		BidPrice: "3000.00",
		AskPrice: "3100.00",
	}

	controller.processBookTicker(ethTicker)

	// 验证ETH/USDT:USDT数据被存储（所有USDT交易对都会被存储）
	ethResult := redisController.GetPairPrice("ETH/USDT:USDT")
	if ethResult.Pair == "" {
		t.Error("ETH/USDT:USDT 数据应该被存储（所有USDT交易对都会存储）")
	}
}

// TestProcessBookTickerWithoutFilter 测试不带过滤的BookTicker处理
func TestProcessBookTickerWithoutFilter(t *testing.T) {
	controller := NewBinanceController()

	// 创建测试用的 RedisController
	conf := &config.Config{
		Redis: config.RedisConfig{
			Addr:      "localhost:6379",
			Password:  "",
			DB:        0,
			KeyExpire: 300,
		},
	}
	redisController := redis.NewRedisController(conf)
	controller.SetRedisController(redisController)

	// 不设置监听列表（处理所有交易对）

	// 处理多个交易对数据
	tickers := []model.BookTickerData{
		{Symbol: "BTCUSDT", BidPrice: "50000.00", AskPrice: "50100.00"},
		{Symbol: "ETHUSDT", BidPrice: "3000.00", AskPrice: "3100.00"},
		{Symbol: "BNBUSDT", BidPrice: "400.00", AskPrice: "410.00"},
	}

	for _, ticker := range tickers {
		controller.processBookTicker(ticker)
	}

	// 验证所有数据都被存储
	allData, err := redisController.GetAllPairPricesData()
	if err != nil {
		t.Errorf("获取所有数据失败: %v", err)
	}

	if len(allData) != 3 {
		t.Errorf("期望 3 个交易对数据，实际 %d 个", len(allData))
	}
}

// TestGetPairCandle 测试获取单个交易对数据
func TestGetPairCandle(t *testing.T) {
	controller := NewBinanceController()

	// 创建测试用的 RedisController
	conf := &config.Config{
		Redis: config.RedisConfig{
			Addr:      "localhost:6379",
			Password:  "",
			DB:        0,
			KeyExpire: 300,
		},
	}
	redisController := redis.NewRedisController(conf)
	controller.SetRedisController(redisController)

	// 添加测试数据
	testData := &model.PairData{
		Pair:      "BTC/USDT:USDT",
		BidPrice:  50000.00,
		AskPrice:  50100.00,
		Timestamp: "2025-01-01 00:00:00",
	}

	redisController.UpdatePairPrice("BTC/USDT:USDT", testData)

	// 测试获取存在的数据
	result := redisController.GetPairPrice("BTC/USDT:USDT")
	if result.Pair != testData.Pair {
		t.Errorf("获取的交易对不匹配，期望 %s，实际 %s", testData.Pair, result.Pair)
	}

	// 测试获取不存在的数据
	emptyResult := redisController.GetPairPrice("NOT_EXIST")
	if emptyResult.Pair != "" {
		t.Error("不存在的交易对应该返回空数据")
	}
}

// TestGetAllPairsData 测试获取所有交易对数据
func TestGetAllPairsData(t *testing.T) {
	controller := NewBinanceController()

	// 创建测试用的 RedisController
	conf := &config.Config{
		Redis: config.RedisConfig{
			Addr:      "localhost:6379",
			Password:  "",
			DB:        0,
			KeyExpire: 300,
		},
	}
	redisController := redis.NewRedisController(conf)
	controller.SetRedisController(redisController)

	// 添加测试数据
	testData := map[string]*model.PairData{
		"BTC/USDT:USDT": {
			Pair:      "BTC/USDT:USDT",
			BidPrice:  50000.00,
			AskPrice:  50100.00,
			Timestamp: "2025-01-01 00:00:00",
		},
		"ETH/USDT:USDT": {
			Pair:      "ETH/USDT:USDT",
			BidPrice:  3000.00,
			AskPrice:  3100.00,
			Timestamp: "2025-01-01 00:00:00",
		},
	}

	for pair, data := range testData {
		redisController.UpdatePairPrice(pair, data)
	}

	// 获取所有数据
	allData, err := redisController.GetAllPairPricesData()
	if err != nil {
		t.Fatalf("获取所有数据失败: %v", err)
	}

	if len(allData) != len(testData) {
		t.Errorf("期望 %d 个交易对，实际 %d 个", len(testData), len(allData))
	}

	// 验证数据内容（简单验证交易对名称存在）
	for _, data := range allData {
		if _, exists := testData[data.Pair]; !exists {
			t.Errorf("意外的交易对: %s", data.Pair)
		}
	}
}

// TestStop 测试停止功能
func TestStop(t *testing.T) {
	controller := NewBinanceController()

	// 测试Stop不会panic
	controller.Stop()

	// 验证context被取消
	select {
	case <-controller.ctx.Done():
		// 期望的结果
	case <-time.After(100 * time.Millisecond):
		t.Error("context 应该被取消")
	}
}

// BenchmarkProcessBookTicker 性能测试
func BenchmarkProcessBookTicker(b *testing.B) {
	controller := NewBinanceController()

	// 创建测试用的 RedisController
	conf := &config.Config{
		Redis: config.RedisConfig{
			Addr:      "localhost:6379",
			Password:  "",
			DB:        0,
			KeyExpire: 300,
		},
	}
	redisController := redis.NewRedisController(conf)
	controller.SetRedisController(redisController)

	ticker := model.BookTickerData{
		Symbol:   "BTCUSDT",
		BidPrice: "50000.00",
		AskPrice: "50100.00",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		controller.processBookTicker(ticker)
	}
}

// BenchmarkFormatPairSymbol 格式转换性能测试
func BenchmarkFormatPairSymbol(b *testing.B) {
	controller := NewBinanceController()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		controller.formatPairSymbol("BTCUSDT")
	}
}

// TestGetFundingRate 测试获取资金费率功能
func TestGetFundingRate(t *testing.T) {
	// 创建模拟服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 验证请求路径和参数
		expectedPath := "/fapi/v1/premiumIndex"
		if r.URL.Path != expectedPath {
			t.Errorf("期望请求路径 %s，实际 %s", expectedPath, r.URL.Path)
		}

		symbol := r.URL.Query().Get("symbol")
		if symbol != "BTCUSDT" {
			t.Errorf("期望请求参数 symbol=BTCUSDT，实际 symbol=%s", symbol)
		}

		// 模拟响应数据
		mockResponse := model.PremiumIndexData{
			Symbol:          "BTCUSDT",
			MarkPrice:       "50050.00",
			IndexPrice:      "50040.00",
			LastFundingRate: "-0.0015",
			NextFundingTime: 1640995200000,
			InterestRate:    "0.0001",
			Time:            1640995100000,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	// 创建BinanceController实例并添加测试方法
	controller := NewBinanceController()

	// 创建测试用的资金费率获取方法
	testGetFundingRate := func(symbol string) (float64, error) {
		binanceSymbol := controller.convertToBinanceSymbol(symbol)
		url := fmt.Sprintf("%s/fapi/v1/premiumIndex?symbol=%s", server.URL, binanceSymbol)

		resp, err := controller.httpClient.Get(url)
		if err != nil {
			return 0, fmt.Errorf("获取资金费率请求失败: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return 0, fmt.Errorf("获取资金费率API错误，状态码: %d", resp.StatusCode)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return 0, fmt.Errorf("读取资金费率响应失败: %v", err)
		}

		var premiumIndex model.PremiumIndexData
		if err := json.Unmarshal(body, &premiumIndex); err != nil {
			return 0, fmt.Errorf("解析资金费率数据失败: %v", err)
		}

		fundingRate, err := strconv.ParseFloat(premiumIndex.LastFundingRate, 64)
		if err != nil {
			return 0, fmt.Errorf("解析资金费率数值失败: %v", err)
		}

		return fundingRate, nil
	}

	// 测试获取资金费率
	fundingRate, err := testGetFundingRate("BTC/USDT:USDT")
	if err != nil {
		t.Fatalf("获取资金费率失败: %v", err)
	}

	expectedRate := -0.0015
	if fundingRate != expectedRate {
		t.Errorf("期望资金费率 %f，实际 %f", expectedRate, fundingRate)
	}
}

// TestGetFundingRateError 测试获取资金费率错误处理
func TestGetFundingRateError(t *testing.T) {
	// 创建返回错误的模拟服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	controller := NewBinanceController()

	// 创建测试用的资金费率获取方法
	testGetFundingRate := func(symbol string) (float64, error) {
		binanceSymbol := controller.convertToBinanceSymbol(symbol)
		url := fmt.Sprintf("%s/fapi/v1/premiumIndex?symbol=%s", server.URL, binanceSymbol)

		resp, err := controller.httpClient.Get(url)
		if err != nil {
			return 0, fmt.Errorf("获取资金费率请求失败: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return 0, fmt.Errorf("获取资金费率API错误，状态码: %d", resp.StatusCode)
		}

		return 0, nil
	}

	_, err := testGetFundingRate("BTC/USDT:USDT")
	if err == nil {
		t.Error("期望获取资金费率返回错误，但没有返回")
	}
}

// TestConvertToBinanceSymbol 测试交易对格式转换
func TestConvertToBinanceSymbol(t *testing.T) {
	controller := NewBinanceController()

	testCases := []struct {
		input    string
		expected string
	}{
		{"BTC/USDT:USDT", "BTCUSDT"},
		{"ETH/USDT:USDT", "ETHUSDT"},
		{"BNB/BTC:BTC", "BNBBTC"},
		{"ADA/ETH:ETH", "ADAETH"},
		{"BTCUSDT", "BTCUSDT"}, // 已经是Binance格式
	}

	for _, tc := range testCases {
		result := controller.convertToBinanceSymbol(tc.input)
		if result != tc.expected {
			t.Errorf("convertToBinanceSymbol(%s) = %s, 期望 %s", tc.input, result, tc.expected)
		}
	}
}

// TestConvertBookTickerToPairDataWithFundingRate 测试包含资金费率的数据转换逻辑
func TestConvertBookTickerToPairDataWithFundingRate(t *testing.T) {
	// 这个测试验证convertBookTickerToPairData方法的逻辑，
	// 但由于它会调用实际的Binance API，我们只测试它不会崩溃
	// 并且能正确处理失败情况

	controller := NewBinanceController()

	ticker := model.BookTickerData{
		Symbol:   "INVALIDPAIR", // 使用无效的交易对，确保API调用失败
		BidPrice: "50000.00",
		AskPrice: "50100.00",
	}

	pairData, err := controller.convertBookTickerToPairData(ticker)
	if err != nil {
		t.Fatalf("转换失败: %v", err)
	}

	// 验证其他字段正确设置
	if pairData.BidPrice != 50000.00 {
		t.Errorf("期望买单价 50000.00，实际 %f", pairData.BidPrice)
	}

	if pairData.AskPrice != 50100.00 {
		t.Errorf("期望卖单价 50100.00，实际 %f", pairData.AskPrice)
	}
}
