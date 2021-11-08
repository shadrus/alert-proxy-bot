package senders

import (
	"errors"
	"fmt"
	"os"
)

const supportedPlatforms = "telegram"

type AlertSender interface {
	Send (message string, chatId string) error
}

func NewAlertSender(target string) (AlertSender, error){
	switch target {
	case "telegram":
		token := os.Getenv("TELEGRAM_TOKEN")
		if token == ""{
			return nil, errors.New("environment variable TELEGRAM_TOKEN was not found")
		}
		return telegramSender{token: token}, nil
	default:
		return nil, errors.New(fmt.Sprintf("We supports only next platforms: %s", supportedPlatforms))
	}
}