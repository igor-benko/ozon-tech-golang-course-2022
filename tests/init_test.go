package tests

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/jackc/pgx/v4/pgxpool"
	"gitlab.ozon.dev/igor.benko.1991/homework/internal/config"
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

	// Инициализация сервиса
	ps := service.NewPersonService(personRepo, vehicleRepo, *cfg)
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
