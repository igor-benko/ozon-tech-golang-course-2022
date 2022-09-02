package bot

import (
	"context"
	"errors"
	"expvar"
	"fmt"
	"io"
	"net/http"
	"time"

	"gitlab.ozon.dev/igor.benko.1991/homework/internal/commander"
	"gitlab.ozon.dev/igor.benko.1991/homework/internal/config"
	"gitlab.ozon.dev/igor.benko.1991/homework/internal/pkg/broker/kafka"
	"gitlab.ozon.dev/igor.benko.1991/homework/internal/pkg/client"
	"google.golang.org/grpc"

	"github.com/Shopify/sarama"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/opentracing/opentracing-go"
	pb "gitlab.ozon.dev/igor.benko.1991/homework/pkg/api"
	"gitlab.ozon.dev/igor.benko.1991/homework/pkg/logger"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"

	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
)

func Run(ctx context.Context, cfg config.Config) {
	bot, err := tgbotapi.NewBotAPI(cfg.Telegram.ApiKey)
	if err != nil {
		logger.FatalKV(err.Error())
	}
	bot.Debug = true

	closer, err := initTracing(cfg)
	if err != nil {
		logger.FatalKV(err.Error())
	}

	defer closer.Close()

	// Инициализация сервисов
	retry_opts := []grpc_retry.CallOption{
		grpc_retry.WithBackoff(grpc_retry.BackoffLinear(time.Duration(cfg.Telegram.RetryIntervalMs) * time.Millisecond)),
		grpc_retry.WithMax(cfg.Telegram.RetryMax),
	}

	tracing_opts := []grpc_opentracing.Option{
		grpc_opentracing.WithTracer(opentracing.GlobalTracer()),
	}

	conn, err := grpc.Dial(cfg.Telegram.PersonService,
		grpc.WithInsecure(),
		grpc.WithStreamInterceptor(grpc_middleware.ChainStreamClient(
			grpc_retry.StreamClientInterceptor(retry_opts...),
			grpc_opentracing.StreamClientInterceptor(tracing_opts...),
		)),

		grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(
			grpc_retry.UnaryClientInterceptor(retry_opts...),
			grpc_opentracing.UnaryClientInterceptor(tracing_opts...),
		)))
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
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	srv := createExpvarServer(cfg.Telegram.ExpvarPort)
	go func() {
		defer cancel()
		logger.Infof("Grpc gateway started!")
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Errorf(err.Error())
		}
	}()

	logger.Infof("Bot started!")
	// Так называемый, gracefull shutdown

	<-ctx.Done()
	logger.Infof("ctx.Done:\n")

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

func initTracing(cfg config.Config) (io.Closer, error) {
	c := &jaegercfg.Configuration{
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans: true,
		},
	}

	return c.InitGlobalTracer(cfg.Telegram.AppName)
}
