package bot

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"gitlab.ozon.dev/igor.benko.1991/homework/internal/commander"
	"gitlab.ozon.dev/igor.benko.1991/homework/internal/config"
	repo "gitlab.ozon.dev/igor.benko.1991/homework/internal/repository/memory"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Run(cfg config.Config) {
	bot, err := tgbotapi.NewBotAPI(cfg.Telegram.ApiKey)
	if err != nil {
		log.Fatal(err)
	}
	bot.Debug = true

	// Инициализация хранилища
	memoryStorage := repo.NewPersonRepo(cfg.Storage)

	// Инициализация обработчиков команд
	personHandler := commander.NewPersonHandler(memoryStorage)

	// Запуск бота
	commander := commander.NewCommander(bot, personHandler, cfg)
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
