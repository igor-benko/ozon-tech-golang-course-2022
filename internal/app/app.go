package app

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"gitlab.ozon.dev/igor.benko.1991/homework/internal/commander"
	"gitlab.ozon.dev/igor.benko.1991/homework/internal/config"
	"gitlab.ozon.dev/igor.benko.1991/homework/internal/handlers"
	"gitlab.ozon.dev/igor.benko.1991/homework/internal/service"
	"gitlab.ozon.dev/igor.benko.1991/homework/internal/storage"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Run() {
	cfg := config.Get()

	bot, err := tgbotapi.NewBotAPI(cfg.Telegram.ApiKey)
	if err != nil {
		log.Fatal(err)
	}
	bot.Debug = true

	// Инициализация хранилища
	memoryStorage := storage.NewMemoryStorage()

	// Инициализация сервисов
	personService := service.NewPersonService(memoryStorage)

	// Инициализация обработчиков команд
	personHandler := handlers.NewPersonHandler(bot, personService)

	// Запуск бота
	commander := commander.NewCommander(bot, personHandler)
	go commander.Run()

	log.Println("Bot started!")
	// Так называемый, gracefull shutdown
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	s := <-interrupt
	log.Println("App received signal: " + s.String())

	// Закрываем канал с updates
	commander.Stop()
}
