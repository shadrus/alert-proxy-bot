package senders

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"strconv"
)

type telegramSender struct {
	token string
}

func (s telegramSender) Send(message string, chatId string) error  {
	chat, err := strconv.ParseInt(chatId, 10, 64)
	if err != nil {
		return err
	}
	bot, err := tgbotapi.NewBotAPI(s.token)
	if err != nil {
		return err
	}
	msg := tgbotapi.NewMessage(chat,message)
	_, err = bot.Send(msg)
	return err
}
