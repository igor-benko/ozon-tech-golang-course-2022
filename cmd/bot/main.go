package main

import (
	"gitlab.ozon.dev/igor.benko.1991/homework/internal/app/bot"
	"gitlab.ozon.dev/igor.benko.1991/homework/internal/config"
	"gitlab.ozon.dev/igor.benko.1991/homework/pkg/logger"
)

func main() {
	cfg, err := config.Init()
	if err != nil {
		logger.FatalKV(err.Error())
	}

	logger.WithAppName(cfg.Telegram.AppName)

	bot.Run(*cfg)
}
