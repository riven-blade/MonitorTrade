package freqtrade

import (
	"encoding/json"
	"monitor-trade/model"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestGetWhitelist 测试获取whitelist功能
func TestGetWhitelist(t *testing.T) {
	// 创建模拟服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 验证请求路径
		if r.URL.Path != "/api/v1/whitelist" {
			t.Errorf("期望请求路径 /api/v1/whitelist，实际 %s", r.URL.Path)
		}

		// 验证请求方法
		if r.Method != "GET" {
			t.Errorf("期望GET请求，实际 %s", r.Method)
		}

		// 验证Authorization头
		authHeader := r.Header.Get("Authorization")
		if authHeader != "Bearer test-access-token" {
			t.Errorf("期望Authorization头包含Bearer token")
		}

		// 模拟响应数据
		mockResponse := model.WhitelistResponse{
			Whitelist: []string{
				"BTC/USDT:USDT",
				"ETH/USDT:USDT",
				"SOL/USDT:USDT",
				"XRP/USDT:USDT",
				"DOGE/USDT:USDT",
			},
			Length: 5,
			Method: []string{
				"VolumePairList",
				"FullTradesFilter",
				"AgeFilter",
				"PriceFilter",
				"SpreadFilter",
			},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	// 创建FreqtradeController实例
	fc := NewFreqtradeController(server.URL, "testuser", "testpass", nil)
	fc.AccessToken = "test-access-token" // 设置测试token

	// 测试GetWhitelist函数
	whitelist, err := fc.getWhitelist()
	if err != nil {
		t.Fatalf("GetWhitelist失败: %v", err)
	}

	// 验证返回的whitelist
	expectedLength := 5
	if len(whitelist) != expectedLength {
		t.Errorf("期望whitelist长度 %d，实际 %d", expectedLength, len(whitelist))
	}

	// 验证具体的交易对
	expectedPairs := []string{
		"BTC/USDT:USDT",
		"ETH/USDT:USDT",
		"SOL/USDT:USDT",
		"XRP/USDT:USDT",
		"DOGE/USDT:USDT",
	}

	for i, expectedPair := range expectedPairs {
		if i >= len(whitelist) {
			t.Errorf("whitelist缺少交易对: %s", expectedPair)
			continue
		}
		if whitelist[i] != expectedPair {
			t.Errorf("交易对 %d: 期望 %s，实际 %s", i, expectedPair, whitelist[i])
		}
	}
}

// TestGetWhitelistError 测试GetWhitelist错误处理
func TestGetWhitelistError(t *testing.T) {
	// 创建返回错误的模拟服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Unauthorized"))
	}))
	defer server.Close()

	// 创建FreqtradeController实例
	fc := NewFreqtradeController(server.URL, "testuser", "testpass", nil)
	fc.AccessToken = "invalid-token"

	// 测试GetWhitelist函数应该返回错误
	_, err := fc.getWhitelist()
	if err == nil {
		t.Error("期望GetWhitelist返回错误，但没有返回")
	}
}

// TestGetWhitelistInvalidJSON 测试无效JSON响应处理
func TestGetWhitelistInvalidJSON(t *testing.T) {
	// 创建返回无效JSON的模拟服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("invalid json"))
	}))
	defer server.Close()

	// 创建FreqtradeController实例
	fc := NewFreqtradeController(server.URL, "testuser", "testpass", nil)
	fc.AccessToken = "test-access-token"

	// 测试GetWhitelist函数应该返回解析错误
	_, err := fc.getWhitelist()
	if err == nil {
		t.Error("期望GetWhitelist返回JSON解析错误，但没有返回")
	}
}
