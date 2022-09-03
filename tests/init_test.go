package tests

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/go-redis/redis/v9"
	"github.com/jackc/pgx/v4/pgxpool"
	"gitlab.ozon.dev/igor.benko.1991/homework/internal/config"
	"gitlab.ozon.dev/igor.benko.1991/homework/internal/pkg/cache"
	repo "gitlab.ozon.dev/igor.benko.1991/homework/internal/repository"
	memory_repo "gitlab.ozon.dev/igor.benko.1991/homework/internal/repository/memory"
	postgres_repo "gitlab.ozon.dev/igor.benko.1991/homework/internal/repository/postgres"
	"gitlab.ozon.dev/igor.benko.1991/homework/internal/service"
	pb "gitlab.ozon.dev/igor.benko.1991/homework/pkg/api"
	"gitlab.ozon.dev/igor.benko.1991/homework/pkg/logger"
	"gitlab.ozon.dev/igor.benko.1991/homework/pkg/postgres"
	"google.golang.org/grpc"
)

var (
	m              sync.Mutex
	pool           *pgxpool.Pool
	personRepo     repo.PersonRepo
	vehicleRepo    repo.VehicleRepo
	personService  pb.PersonServiceServer
	vehicleService pb.VehicleServiceServer
)

func init() {
	cfg, err := config.Init()
	if err != nil {
		logger.FatalKV(err.Error())
	}

	if cfg.PersonService.Storage == config.StoragePostgres {
		err := postgres.Migrate(context.Background(), &cfg.Database)
		if err != nil {
			logger.FatalKV(err.Error())
		}

		pool, err = postgres.New(context.Background(), &cfg.Pooler)
		if err != nil {
			logger.FatalKV(err.Error())
		}

		personRepo = postgres_repo.NewPersonRepo(pool)
		vehicleRepo = postgres_repo.NewVehicleRepo(pool)

	} else if cfg.PersonService.Storage == config.StorageMemory {
		personRepo = memory_repo.NewPersonRepo(cfg.Storage)
		vehicleRepo = memory_repo.NewVehicleRepo(cfg.Storage)
	} else {
		log.Fatalf("Unsupported storage type %s", cfg.PersonService.Storage)
	}

	cache, _ := initCache(*cfg)

	// Инициализация сервиса
	ps := service.NewPersonService(personRepo, vehicleRepo, cache, *cfg)
	vs := service.NewVehicleService(vehicleRepo, *cfg)

	personService = &ps
	vehicleService = &vs

	// GRPC
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.PersonService.Port))
	if err != nil {
		logger.FatalKV(err.Error())
	}

	grpcServer := grpc.NewServer()
	pb.RegisterPersonServiceServer(grpcServer, personService)
	pb.RegisterVehicleServiceServer(grpcServer, vehicleService)

	go func() {
		logger.Infof("Grpc started!")
		if err = grpcServer.Serve(listener); err != nil {
			logger.Errorf(err.Error())
		}
	}()
}

func setUpPool() {
	m.Lock()
}

func tearDownPool() {
	if _, err := pool.Exec(context.TODO(), "DELETE FROM vehicles"); err != nil {
		logger.FatalKV(err.Error())
	}

	if _, err := pool.Exec(context.TODO(), "DELETE FROM persons"); err != nil {
		logger.FatalKV(err.Error())
	}

	m.Unlock()
}

type CloserWithErr func() error

func initCache(cfg config.Config) (cache.CacheClient, CloserWithErr) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Cache.Host, cfg.Cache.Port),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return cache.NewRedisCache(cfg, client), client.Close
}
