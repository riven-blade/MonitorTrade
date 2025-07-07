package tg

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (tg *TgController) SendWelcomeWithButtons() error {
	msg := tgbotapi.NewMessage(tg.TgId, "Starting Trade Monitor!")
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("/show"),
			tgbotapi.NewKeyboardButton("/whitelist"),
			tgbotapi.NewKeyboardButton("/adjust"),
		),
	)
	msg.ReplyMarkup = keyboard
	_, err := tg.Bot.Send(msg)
	return err
}

func (tg *TgController) SendMessage(msg string) {
	tgMsg := tgbotapi.NewMessage(tg.TgId, msg)
	if _, err := tg.Bot.Send(tgMsg); err != nil {
		log.Printf("发送消息到 Telegram 失败: %v", err)
	}
}

func (tg *TgController) SendMessageByChan(messageChan chan string) {
	for message := range messageChan {
		tg.SendMessage(message)
	}
}
