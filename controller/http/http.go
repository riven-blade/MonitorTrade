package http

import (
	"encoding/json"
	"log"
	"monitor-trade/controller"
	"monitor-trade/controller/freqtrade"
	"monitor-trade/controller/redis"
	"monitor-trade/model"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func ListenAndServe(hh *HttpHandler) {
	r := gin.Default()
	r.GET("/api/monitor", hh.ListMonitor)
	r.POST("/api/webhook", hh.HandleWebhook) // 单一webhook端点

	s := &http.Server{
		Addr:           ":8888",
		Handler:        r,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()
}

type HttpHandler struct {
	MainController  *controller.MainController
	redisController *redis.RedisController
	fc              *freqtrade.FreqtradeController
}

func NewHttpHandler(mc *controller.MainController, redisController *redis.RedisController, fc *freqtrade.FreqtradeController) *HttpHandler {
	return &HttpHandler{
		MainController:  mc,
		redisController: redisController,
		fc:              fc,
	}
}

func (h *HttpHandler) ListMonitor(c *gin.Context) {
	ctx := c.Request.Context()
	prefix := c.DefaultQuery("prefix", "*")
	keys, err := h.MainController.RedisController.Client.Keys(ctx, redis.MonitorKey+prefix).Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var pairMonitorDataList []model.PairMonitorDataDetail
	for i := range keys {
		var pairMonitorData model.PairMonitorData
		key := keys[i]
		data, err2 := h.MainController.RedisController.Client.Get(ctx, key).Result()
		if err2 != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err2.Error()})
			return
		}
		if data == "" {
			continue
		}
		json.Unmarshal([]byte(data), &pairMonitorData)

		pairMonitorTTL, _ := h.MainController.RedisController.Client.TTL(ctx, key).Result()
		pairData := h.redisController.GetPairPrice(pairMonitorData.Pair)
		pairMonitorDataList = append(pairMonitorDataList, model.PairMonitorDataDetail{
			PairData:        pairData,
			PairMonitorData: pairMonitorData,
			TTL:             pairMonitorTTL.Seconds(),
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": pairMonitorDataList})
}

// HandleWebhook 处理Freqtrade webhook消息
func (h *HttpHandler) HandleWebhook(c *gin.Context) {
	var rawData map[string]interface{}
	if err := c.ShouldBindJSON(&rawData); err != nil {
		log.Printf("解析webhook数据失败: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	// 打印原始数据
	for s := range rawData {
		log.Println(rawData[s])
	}
	go h.fc.CheckRedisPairStatus()
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
