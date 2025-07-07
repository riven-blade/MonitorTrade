package freqtrade

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"monitor-trade/controller/redis"
	"monitor-trade/model"
	"net/http"
	"time"
)

type FreqtradeController struct {
	BaseUrl         string
	Username        string
	Password        string
	AccessToken     string
	RefreshToken    string
	stopChan        chan struct{}
	stopChanPair    chan struct{}
	httpClient      *http.Client
	PositionStatus  model.PositionStatus
	TradeStatus     []model.TradePosition
	redisController *redis.RedisController
	messageChan     chan string
}

func NewFreqtradeController(baseUrl, username, password string, redisController *redis.RedisController) *FreqtradeController {
	return &FreqtradeController{
		BaseUrl:         baseUrl,
		Username:        username,
		Password:        password,
		stopChan:        make(chan struct{}),
		redisController: redisController,
		httpClient:      &http.Client{Timeout: 10 * time.Second},
	}
}

func (fc *FreqtradeController) startTokenRefresher() {
	if fc.stopChan != nil {
		close(fc.stopChan) // 防止重复启动
	}
	fc.stopChan = make(chan struct{})

	go func() {
		log.Println("Token 刷新器已启动")
		ticker := time.NewTicker(10 * time.Minute)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				go fc.refreshToken()
			case <-fc.stopChan:
				log.Println("Token 刷新器已停止")
				return
			}
		}
	}()
}

func (fc *FreqtradeController) pairRefresher() {
	if fc.stopChanPair != nil {
		close(fc.stopChanPair) // 防止重复启动
	}
	fc.stopChanPair = make(chan struct{})

	go func() {
		log.Println("交易对刷新器已启动")
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				go fc.setPairWhiteList()
			case <-fc.stopChanPair:
				log.Println("交易对刷新器已停止")
				return
			}
		}
	}()
}

func (fc *FreqtradeController) doRequest(method, url string, body io.Reader, useAccessToken bool) ([]byte, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	if useAccessToken {
		req.Header.Set("Authorization", "Bearer "+fc.AccessToken)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := fc.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s %s 请求失败: %s", method, url, string(respBody))
	}
	return respBody, nil
}

func (fc *FreqtradeController) Init(messageChan chan string) {
	fc.messageChan = messageChan
	url := fmt.Sprintf("%v/api/v1/token/login", fc.BaseUrl)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.SetBasicAuth(fc.Username, fc.Password)

	resp, err := fc.httpClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("登录失败: %d %s", resp.StatusCode, resp.Status)
	}

	body, _ := io.ReadAll(resp.Body)
	var loginResp model.LoginResponse
	if err := json.Unmarshal(body, &loginResp); err != nil {
		log.Fatalf("解析登录响应失败: %v", err)
	}

	fc.AccessToken = loginResp.AccessToken
	fc.RefreshToken = loginResp.RefreshToken

	log.Println("首次登录成功")

	// 启动交易对刷新器和token刷新器
	go fc.CheckRedisPairStatus()
	go fc.setPairWhiteList()
	go fc.pairRefresher()
	go fc.startTokenRefresher()
}

func (fc *FreqtradeController) refreshToken() {
	url := fmt.Sprintf("%v/api/v1/token/refresh", fc.BaseUrl)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		log.Printf("创建刷新请求失败: %v", err)
		return
	}
	req.Header.Set("Authorization", "Bearer "+fc.RefreshToken)

	resp, err := fc.httpClient.Do(req)
	if err != nil {
		log.Printf("刷新 token 请求失败: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("刷新 token 失败: %v", resp.Status)
		return
	}

	body, _ := io.ReadAll(resp.Body)
	var loginResp model.LoginResponse
	if err := json.Unmarshal(body, &loginResp); err != nil {
		log.Printf("解析刷新响应失败: %v", err)
		return
	}

	fc.AccessToken = loginResp.AccessToken
	log.Println("刷新 token 成功")
}

func (fc *FreqtradeController) ForceBuy(pair string, price float64, side string) error {
	url := fmt.Sprintf("%s/api/v1/forcebuy", fc.BaseUrl)
	payload := model.ForceBuyPayload{
		Pair:      pair,
		Price:     price,
		OrderType: "limit",
		Side:      side,
		EntryTag:  "force_entry",
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	respBody, err := fc.doRequest("POST", url, bytes.NewReader(body), true)
	if err != nil {
		return err
	}

	log.Printf("forcebuy 成功: %s", string(respBody))
	return nil
}

func (fc *FreqtradeController) ForceAdjustBuy(pair string, price float64, side string, stakeAmount float64) error {
	url := fmt.Sprintf("%s/api/v1/forcebuy", fc.BaseUrl)
	payload := model.ForceAdjustBuyPayload{
		Pair:        pair,
		Price:       price,
		OrderType:   "limit",
		Side:        side,
		EntryTag:    "force_entry",
		StakeAmount: stakeAmount,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	respBody, err := fc.doRequest("POST", url, bytes.NewReader(body), true)
	if err != nil {
		return err
	}

	log.Printf("forceadjustbuy 成功: %s", string(respBody))
	return nil
}

func (fc *FreqtradeController) ForceSell(tradeId string, orderType string, amount string) error {
	url := fmt.Sprintf("%s/api/v1/forcesell", fc.BaseUrl)
	payload := model.ForceSellPayload{
		TradeId:   tradeId,
		OrderType: orderType,
		Amount:    amount,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	respBody, err := fc.doRequest("POST", url, bytes.NewReader(body), true)
	if err != nil {
		return err
	}

	log.Printf("forcesell 成功: %s", string(respBody))
	return nil
}

func (fc *FreqtradeController) getCount() error {
	url := fmt.Sprintf("%v/api/v1/count", fc.BaseUrl)
	body, err := fc.doRequest("GET", url, nil, true)
	if err != nil {
		return err
	}

	var positions model.PositionStatus
	if err = json.Unmarshal(body, &positions); err != nil {
		return err
	}
	fc.PositionStatus = positions
	return nil
}

func (fc *FreqtradeController) getStatus() error {
	url := fmt.Sprintf("%s/api/v1/status", fc.BaseUrl)
	body, err := fc.doRequest("GET", url, nil, true)
	if err != nil {
		return err
	}

	var trades []model.TradePosition
	if err := json.Unmarshal(body, &trades); err != nil {
		return err
	}
	fc.TradeStatus = trades
	return nil
}

func (fc *FreqtradeController) fetchTradeData() error {
	err := fc.getStatus()
	if err != nil {
		return err
	}
	// 获取当前持仓数量
	err = fc.getCount()
	if err != nil {
		return err
	}
	return nil
}

func (fc *FreqtradeController) CheckRedisPairStatus() {
	err := fc.fetchTradeData()
	if err != nil {
		log.Printf("获取交易数据失败: %v", err)
	}

	tradeStatus := fc.TradeStatus
	if tradeStatus == nil {
		log.Println("获取交易状态失败，无法检查Redis交易对状态")
		return
	}
	// 遍历当前交易状态，检查是否有需要更新的交易对
	for i := range tradeStatus {
		trade := tradeStatus[i]
		if len(trade.Orders) >= 1 {
			if !trade.Orders[0].IsOpen {
				if trade.IsShort {
					// short 交易对
					_, exits := fc.redisController.GetMonitorPair(trade.Pair, "short")
					if exits {
						log.Printf("交易对 %s 的做空仓位已经成交，删除 Redis 中的监控数据", trade.Pair)
						fc.redisController.DeleteMonitorPair(trade.Pair, "short")
						go func() {
							// 发送做空成功消息
							fc.messageChan <- fmt.Sprintf("✅ %s 做空仓位已成交，删除 Redis 中的监控数据", trade.Pair)
						}()
					}
				} else {
					// long 交易对
					_, exits := fc.redisController.GetMonitorPair(trade.Pair, "long")
					if exits {
						log.Printf("交易对 %s 的做多仓位已经成交，删除 Redis 中的监控数据", trade.Pair)
						fc.redisController.DeleteMonitorPair(trade.Pair, "long")
						go func() {
							// 发送做多成功消息
							fc.messageChan <- fmt.Sprintf("✅ %s 做多仓位已成交，删除 Redis 中的监控数据", trade.Pair)
						}()
					}
				}
			}
		}
	}
}

// 检查是否可以强制买入
func (fc *FreqtradeController) CheckForceBuy(pair string) bool {
	err := fc.fetchTradeData()
	if err != nil {
		log.Printf("获取交易数据失败: %v", err)
		return false
	}

	tradeStatus := fc.TradeStatus
	for i := range tradeStatus {
		trade := tradeStatus[i]
		if trade.Pair == pair {
			return false
		}
	}

	return len(tradeStatus) < fc.PositionStatus.Max
}

// GetWhitelist 获取交易对白名单
func (fc *FreqtradeController) getWhitelist() ([]string, error) {
	url := fmt.Sprintf("%s/api/v1/whitelist", fc.BaseUrl)
	body, err := fc.doRequest("GET", url, nil, true)
	if err != nil {
		log.Printf("获取whitelist失败: %v", err)
		return nil, err
	}

	var whitelistResp model.WhitelistResponse
	if err := json.Unmarshal(body, &whitelistResp); err != nil {
		log.Printf("解析whitelist响应失败: %v", err)
		return nil, err
	}

	log.Printf("获取whitelist成功，共 %d 个交易对", whitelistResp.Length)
	return whitelistResp.Whitelist, nil
}

func (fc *FreqtradeController) setPairWhiteList() {
	whitelist, err := fc.getWhitelist()
	if err != nil {
		log.Printf("获取交易对白名单失败: %v", err)
		return
	}
	fc.redisController.SetWatchedPairs(whitelist)
	log.Println("交易对白名单已刷新")
}
