package main

import (
	"log"

	"gitlab.ozon.dev/igor.benko.1991/homework/internal/app/server"
	"gitlab.ozon.dev/igor.benko.1991/homework/internal/config"
)

func main() {
	cfg, err := config.Init()
	if err != nil {
		log.Fatal(err)
	}

	server.Run(*cfg)
}
