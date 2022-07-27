package main

import (
	"log"

	"gitlab.ozon.dev/igor.benko.1991/homework/internal/app/bot"
	"gitlab.ozon.dev/igor.benko.1991/homework/internal/config"
)

func main() {
	if err := config.Init(); err != nil {
		log.Fatal(err)
	}

	bot.Run()
}
