package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gitlab.ozon.dev/igor.benko.1991/homework/docs"
	"gitlab.ozon.dev/igor.benko.1991/homework/internal/config"
	repo "gitlab.ozon.dev/igor.benko.1991/homework/internal/repository"
	memory_repo "gitlab.ozon.dev/igor.benko.1991/homework/internal/repository/memory"
	postgres_repo "gitlab.ozon.dev/igor.benko.1991/homework/internal/repository/postgres"
	"gitlab.ozon.dev/igor.benko.1991/homework/internal/service"
	pb "gitlab.ozon.dev/igor.benko.1991/homework/pkg/api"
	"gitlab.ozon.dev/igor.benko.1991/homework/pkg/logger"
	"gitlab.ozon.dev/igor.benko.1991/homework/pkg/postgres"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

const (
	shotdownTimeout = 5 * time.Second
)

func Run(cfg config.Config) {
	// Инициализация хранилища
	var personRepo repo.PersonRepo
	var vehicleRepo repo.VehicleRepo

	if cfg.PersonService.Storage == config.StoragePostgres {
		err := postgres.Migrate(context.Background(), &cfg.Database)
		if err != nil {
			logger.FatalKV(err.Error())
		}

		pool, err := postgres.New(context.Background(), &cfg.Pooler)
		if err != nil {
			logger.FatalKV(err.Error())
		}

		defer pool.Close()

		personRepo = postgres_repo.NewPersonRepo(pool)
		vehicleRepo = postgres_repo.NewVehicleRepo(pool)

	} else if cfg.PersonService.Storage == config.StorageMemory {
		personRepo = memory_repo.NewPersonRepo(cfg.Storage)
		vehicleRepo = memory_repo.NewVehicleRepo(cfg.Storage)

	} else {
		logger.Fatalf("Unsupported storage type %s", cfg.PersonService.Storage)
	}

	// Инициализация сервиса
	personService := service.NewPersonService(personRepo, vehicleRepo, cfg)
	vehicleService := service.NewVehicleService(vehicleRepo, cfg)

	// GRPC
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.PersonService.Port))
	if err != nil {
		logger.FatalKV(err.Error())
	}

	grpcServer := grpc.NewServer()
	pb.RegisterPersonServiceServer(grpcServer, &personService)
	pb.RegisterVehicleServiceServer(grpcServer, &vehicleService)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		logger.Infof("Grpc started!")
		if err = grpcServer.Serve(listener); err != nil {
			logger.Errorf(err.Error())
			cancel()
		}
	}()

	// GRPC http gateway
	gatewayServer := createGatewayServer(cfg.PersonService.Port, cfg.PersonService.GatewayPort)
	go func() {
		logger.Infof("Grpc gateway started!")
		if err := gatewayServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Errorf(err.Error())
			cancel()
		}
	}()

	// Так называемый, gracefull shutdown
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case v := <-interrupt:
		log.Printf("signal.Notify: %s\n", v)
	case done := <-ctx.Done():
		log.Printf("ctx.Done: %s\n", done)
	}

	// Даем 5 сек на обработку текущих запросов
	ctx, cancel = context.WithTimeout(context.Background(), shotdownTimeout)
	defer cancel()

	if err := gatewayServer.Shutdown(ctx); err != nil {
		log.Printf("Error on gatewayServer stop: %s\n", err)
	} else {
		logger.Infof("gatewayServer stopped")
	}

	grpcServer.GracefulStop()
	logger.Infof("Grpc server stopped")
}

func createGatewayServer(grpcPort, httpPort int) *http.Server {
	ctx := context.Background()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	if err := pb.RegisterPersonServiceHandlerFromEndpoint(ctx, mux, fmt.Sprintf(":%d", grpcPort), opts); err != nil {
		logger.FatalKV(err.Error())
	}

	if err := pb.RegisterVehicleServiceHandlerFromEndpoint(ctx, mux, fmt.Sprintf(":%d", grpcPort), opts); err != nil {
		logger.FatalKV(err.Error())
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.RemoveExtraSlash = true

	router.Group("v1/*{grpc_gateway}").Any("", gin.WrapH(mux))

	router.StaticFS(docs.SwaggerFileName, http.FS(docs.SwaggerFile))
	router.StaticFileFS(docs.SwaggerFileName, docs.SwaggerFileName, http.FS(docs.SwaggerFile))
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL(fmt.Sprintf("/%s", docs.SwaggerFileName))))

	return &http.Server{
		Addr:    fmt.Sprintf(":%d", httpPort),
		Handler: router,
	}
}
