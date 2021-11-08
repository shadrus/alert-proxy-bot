package main

import (
	"AlertProxyBot/src/internal/alertproxybot"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.TextFormatter{
		DisableColors:  false,
		FullTimestamp:  true,
		DisableSorting: true,
	})
	log.Info("Starting proxy bot")
	alertproxybot.StartServer()
}