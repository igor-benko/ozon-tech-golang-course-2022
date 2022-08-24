package person_consumer

import (
	"context"
	"errors"
	"expvar"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Shopify/sarama"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"gitlab.ozon.dev/igor.benko.1991/homework/internal/config"
	"gitlab.ozon.dev/igor.benko.1991/homework/internal/consumer"
	"gitlab.ozon.dev/igor.benko.1991/homework/internal/pkg/broker/kafka"
	repo "gitlab.ozon.dev/igor.benko.1991/homework/internal/repository"
	memory_repo "gitlab.ozon.dev/igor.benko.1991/homework/internal/repository/memory"
	postgres_repo "gitlab.ozon.dev/igor.benko.1991/homework/internal/repository/postgres"
	"gitlab.ozon.dev/igor.benko.1991/homework/pkg/logger"
	"gitlab.ozon.dev/igor.benko.1991/homework/pkg/postgres"
)

const (
	shotdownTimeout = 5 * time.Second
)

func Run(cfg config.Config) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	closer, err := initTracing(cfg)
	if err != nil {
		logger.FatalKV(err.Error())
	}

	defer closer.Close()

	// Инициализация хранилища
	var personRepo repo.PersonRepo

	if cfg.PersonService.Storage == config.StoragePostgres {
		pool, err := postgres.New(context.Background(), &cfg.Pooler)
		if err != nil {
			logger.FatalKV(err.Error())
		}

		defer pool.Close()

		personRepo = postgres_repo.NewPersonRepo(pool)

	} else if cfg.PersonService.Storage == config.StorageMemory {
		personRepo = memory_repo.NewPersonRepo(cfg.Storage)

	} else {
		logger.Fatalf("Unsupported storage type %s", cfg.PersonService.Storage)
	}

	// Инициализация брокера
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	producer, err := sarama.NewSyncProducer(cfg.Kafka.Brokers, config)
	if err != nil {
		logger.FatalKV(err.Error())
	}

	// Person consumer
	personGroup, err := sarama.NewConsumerGroup(cfg.Kafka.Brokers, cfg.PersonConsumer.GroupName, config)
	if err != nil {
		logger.FatalKV(err.Error())
	}
	personBroker := kafka.NewKafkaBroker(producer, personGroup)
	personConsumer := consumer.NewPersonConsumer(cfg, personBroker, personRepo)
	go func() {
		if err := personConsumer.Consume(ctx, cfg.Kafka.IncomeTopic); err != nil {
			cancel()
		}
	}()

	// Verify consumer
	verifyGroup, err := sarama.NewConsumerGroup(cfg.Kafka.Brokers, cfg.VerifyConsumer.GroupName, config)
	if err != nil {
		logger.FatalKV(err.Error())
	}
	verifyBroker := kafka.NewKafkaBroker(producer, verifyGroup)
	verifyConsumer := consumer.NewVerifyConsumer(cfg, verifyBroker, personRepo)
	go func() {
		if err := verifyConsumer.Consume(ctx, cfg.Kafka.VerifyTopic); err != nil {
			cancel()
		}
	}()

	// Rollbak consumer
	rollbackGroup, err := sarama.NewConsumerGroup(cfg.Kafka.Brokers, cfg.RollbackConsumer.GroupName, config)
	if err != nil {
		logger.FatalKV(err.Error())
	}
	rollbackBroker := kafka.NewKafkaBroker(producer, rollbackGroup)
	rollbackConsumer := consumer.NewRollbackConsumer(cfg, rollbackBroker, personRepo)
	go func() {
		if err := rollbackConsumer.Consume(ctx, cfg.Kafka.ErrorTopic); err != nil {
			cancel()
		}
	}()

	srv := createExpvarServer(cfg.PersonConsumer.ExpvarPort)
	go func() {
		logger.Infof("Expvar server started!")
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Errorf(err.Error())
			cancel()
		}
	}()

	// Так называемый, gracefull shutdown
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case v := <-interrupt:
		logger.Infof("signal.Notify: %s\n", v)
	case done := <-ctx.Done():
		logger.Infof("ctx.Done: %s\n", done)
	}

	// Даем 5 сек на обработку текущих запросов
	logger.Infof("Grpc server stopped")
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

	closer, err := c.InitGlobalTracer(cfg.PersonConsumer.AppName)
	if err != nil {
		return nil, err
	}

	return closer, nil
}
