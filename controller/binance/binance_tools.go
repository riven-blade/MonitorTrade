package binance

import (
	"encoding/json"
	"fmt"
	"io"
	"monitor-trade/model"
	"net/http"
	"strconv"
	"strings"
)

// 获取资金费率
func (b *BinanceController) GetFundingRate(symbol string) (float64, error) {
	// 转换交易对格式：BTC/USDT:USDT -> BTCUSDT
	binanceSymbol := b.convertToBinanceSymbol(symbol)

	url := fmt.Sprintf("https://fapi.binance.com/fapi/v1/premiumIndex?symbol=%s", binanceSymbol)

	resp, err := b.httpClient.Get(url)
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

	return fundingRate * 100, nil
}

// 转换交易对格式：BTC/USDT:USDT -> BTCUSDT
func (b *BinanceController) convertToBinanceSymbol(pair string) string {
	// 移除后缀":USDT"等
	if colonIndex := strings.Index(pair, ":"); colonIndex != -1 {
		pair = pair[:colonIndex]
	}

	// 移除斜杠
	return strings.ReplaceAll(pair, "/", "")
}
