package bot

import (
	"context"
	"errors"
	"expvar"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gitlab.ozon.dev/igor.benko.1991/homework/internal/commander"
	"gitlab.ozon.dev/igor.benko.1991/homework/internal/config"
	"gitlab.ozon.dev/igor.benko.1991/homework/internal/pkg/broker/kafka"
	"gitlab.ozon.dev/igor.benko.1991/homework/internal/pkg/client"
	"google.golang.org/grpc"

	"github.com/Shopify/sarama"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	pb "gitlab.ozon.dev/igor.benko.1991/homework/pkg/api"
	"gitlab.ozon.dev/igor.benko.1991/homework/pkg/logger"
)

func Run(cfg config.Config) {
	bot, err := tgbotapi.NewBotAPI(cfg.Telegram.ApiKey)
	if err != nil {
		logger.FatalKV(err.Error())
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
		logger.FatalKV(err.Error())
	}

	personService := pb.NewPersonServiceClient(conn)

	// Инициализация брокера
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(cfg.Kafka.Brokers, config)
	if err != nil {
		logger.FatalKV(err.Error())
	}

	broker := kafka.NewKafkaBroker(producer, nil)
	personClient := client.NewPersonClient(cfg, personService, broker)

	// Инициализация обработчиков команд
	personHandler := commander.NewPersonHandler(personClient)

	// Запуск бота
	commander := commander.NewCommander(bot, personHandler, cfg)
	go commander.Run()

	// Запуск expvar сервера
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	srv := createExpvarServer(cfg.Telegram.ExpvarPort)
	go func() {
		logger.Infof("Grpc gateway started!")
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Errorf(err.Error())
			cancel()
		}
	}()

	logger.Infof("Bot started!")
	// Так называемый, gracefull shutdown
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case v := <-interrupt:
		logger.Infof("signal.Notify: %s\n", v)
	case done := <-ctx.Done():
		logger.Infof("ctx.Done: %s\n", done)
	}

	// Закрываем канал с updates
	commander.Stop()

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
}

func createExpvarServer(port int) *http.Server {
	mux := http.NewServeMux()
	mux.Handle("/stats", expvar.Handler())

	return &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}
}
