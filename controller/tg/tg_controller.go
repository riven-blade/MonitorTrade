package tg

import (
	"log"
	"monitor-trade/config"
	"monitor-trade/controller/freqtrade"
	"monitor-trade/controller/redis"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	LongDirect  = "long"
	ShortDirect = "short"
)

type TgController struct {
	BotToken            string
	TgId                int64
	Bot                 *tgbotapi.BotAPI
	RedisController     *redis.RedisController
	FreqtradeController *freqtrade.FreqtradeController
	Conf                *config.Config
}

func NewTgController(botToken string, tgId int64, controller *redis.RedisController, freqtradeController *freqtrade.FreqtradeController, conf *config.Config) *TgController {
	// 初始化 Telegram 机器人
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatalf("无法初始化 Telegram 机器人: %v", err)
	}
	log.Printf("已授权账号 %s", bot.Self.UserName)

	return &TgController{
		BotToken:            botToken,
		TgId:                tgId,
		Bot:                 bot,
		RedisController:     controller,
		FreqtradeController: freqtradeController,
		Conf:                conf,
	}
}
