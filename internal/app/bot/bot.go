package bot

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gitlab.ozon.dev/igor.benko.1991/homework/internal/commander"
	"gitlab.ozon.dev/igor.benko.1991/homework/internal/config"
	"google.golang.org/grpc"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	pb "gitlab.ozon.dev/igor.benko.1991/homework/pkg/api"
)

func Run(cfg config.Config) {
	bot, err := tgbotapi.NewBotAPI(cfg.Telegram.ApiKey)
	if err != nil {
		log.Fatal(err)
	}
	bot.Debug = true

	// Инициализация сервисов
	opts := []grpc_retry.CallOption{
		grpc_retry.WithBackoff(grpc_retry.BackoffLinear(time.Duration(cfg.Telegram.RetryIntervalMs) * time.Millisecond)),
		grpc_retry.WithMax(cfg.Telegram.RetryMax),
	}

	conn, err := grpc.Dial(cfg.Telegram.PersonService,
		grpc.WithInsecure(),
		grpc.WithStreamInterceptor(grpc_retry.StreamClientInterceptor(opts...)),
		grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(opts...)))
	if err != nil {
		log.Fatal(err)
	}

	personService := pb.NewPersonServiceClient(conn)

	// Инициализация обработчиков команд
	personHandler := commander.NewPersonHandler(personService)

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
