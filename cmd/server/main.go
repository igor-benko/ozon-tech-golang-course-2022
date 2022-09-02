package main

import (
	"context"
	"os/signal"
	"syscall"

	"gitlab.ozon.dev/igor.benko.1991/homework/internal/app/server"
	"gitlab.ozon.dev/igor.benko.1991/homework/internal/config"
	"gitlab.ozon.dev/igor.benko.1991/homework/pkg/logger"
)

func main() {
	cfg, err := config.Init()
	if err != nil {
		logger.FatalKV(err.Error())
	}

	logger.WithAppName(cfg.PersonService.AppName)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM)
	defer stop()

	server.Run(ctx, *cfg)
}
