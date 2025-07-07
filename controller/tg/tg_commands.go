package tg

import (
	"log"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (tg *TgController) HandleCommand() {
	err := tg.SendWelcomeWithButtons()
	if err != nil {
		log.Printf("发送欢迎消息失败: %v", err)
		return
	}
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := tg.Bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil {
			continue
		}
		if !update.Message.IsCommand() {
			continue
		}
		msg := tgbotapi.NewMessage(tg.TgId, "")
		switch update.Message.Command() {
		case "s", "short":
			args := update.Message.CommandArguments()
			parts := strings.Split(args, " ")
			if len(parts) < 2 {
				msg.Text = "用法: /s [pair] [price]"
			} else {
				pair := tg.HandlePair(parts[0])
				price, err := strconv.ParseFloat(parts[1], 64)
				if err != nil {
					msg.Text = "价格必须是有效的数字"
				} else {
					msg.Text = tg.handleShortCommand(pair, price)
				}
			}
		case "l", "long":
			args := update.Message.CommandArguments()
			parts := strings.Split(args, " ")
			if len(parts) < 2 {
				msg.Text = "用法: /l [pair] [price]"
			} else {
				pair := tg.HandlePair(parts[0])
				price, err := strconv.ParseFloat(parts[1], 64)
				if err != nil {
					msg.Text = "价格必须是有效的数字"
				} else {
					msg.Text = tg.handleLongCommand(pair, price)
				}
			}
		case "c", "cancel":
			args := update.Message.CommandArguments()
			parts := strings.Split(args, " ")
			if len(parts) < 2 {
				msg.Text = "用法: /c [pair] [direction]"
			} else {
				pair := tg.HandlePair(parts[0])
				direct := parts[1]
				msg.Text = tg.handleCancelCommand(pair, direct)
			}
		case "show":
			args := update.Message.CommandArguments()
			parts := strings.Split(args, " ")
			if len(parts) == 0 || (len(parts) < 2 && strings.Trim(parts[0], " ") == "") {
				msg.Text = tg.handleShowConfigCommand()
			} else {
				pair := tg.HandlePair(parts[0])
				msg.Text = tg.handleShowCommand(pair)
			}
		case "adjust":
			// /adjust 现在显示仓位信息
			msg.Text = tg.handleShowPositionsCommand()
		case "whitelist":
			// 处理白名单命令
			msg.Text = tg.handleWhiteList()
		case "ad":
			args := update.Message.CommandArguments()
			parts := strings.Split(args, " ")
			if len(parts) < 2 {
				msg.Text = "用法: /ad [pair] [num] [price]"
			} else {
				pair := tg.HandlePair(parts[0])
				stakeAmount, err1 := strconv.ParseFloat(parts[1], 64)
				if err1 != nil {
					msg.Text = "stakeAmount必须是有效的数字"
				} else {
					var price float64 = 0 // 默认价格为0，表示使用当前市价
					if len(parts) >= 3 {
						// 如果提供了价格参数
						var err2 error
						price, err2 = strconv.ParseFloat(parts[2], 64)
						if err2 != nil {
							msg.Text = "价格必须是有效的数字"
						} else {
							msg.Text = tg.handleADCommand(pair, stakeAmount, price)
						}
					} else {
						// 没有价格参数，使用默认价格0
						msg.Text = tg.handleADCommand(pair, stakeAmount, price)
					}
				}
			}
		case "pc":
			args := update.Message.CommandArguments()
			parts := strings.Split(args, " ")
			if len(parts) < 2 {
				msg.Text = "用法: /pc [pair] [num]"
			} else {
				pair := tg.HandlePair(parts[0])
				stakeAmount, err := strconv.ParseFloat(parts[1], 64)
				if err != nil {
					msg.Text = "stakeAmount必须是有效的数字"
				} else {
					msg.Text = tg.handlePCCommand(pair, stakeAmount)
				}
			}
		default:
			msg.Text = "未知命令。支持的命令: /short /s, /long /l, /cancel /c, /show, /whitelist, /adjust, /ad, /pc"
		}

		log.Println(msg.Text)
		if _, err := tg.Bot.Send(msg); err != nil {
			log.Printf("发送消息失败: %v", err)
		}
	}
}
