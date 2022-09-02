package server

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
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
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/opentracing/opentracing-go"

	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

const (
	shotdownTimeout = 5 * time.Second
)

func Run(ctx context.Context, cfg config.Config) {
	closer, err := initTracing(cfg)
	if err != nil {
		logger.FatalKV(err.Error())
	}

	defer closer.Close()

	// Инициализация хранилища
	repos, close := initRepo(cfg)
	defer close()

	// Инициализация сервиса
	personService := service.NewPersonService(repos.person, repos.vehicle, cfg)
	vehicleService := service.NewVehicleService(repos.vehicle, cfg)

	// GRPC
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.PersonService.Port))
	if err != nil {
		logger.FatalKV(err.Error())
	}

	tracing_opts := []grpc_opentracing.Option{
		grpc_opentracing.WithTracer(opentracing.GlobalTracer()),
	}

	grpcServer := grpc.NewServer(
		grpc.StreamInterceptor(
			grpc_opentracing.StreamServerInterceptor(tracing_opts...),
		),

		grpc.UnaryInterceptor(
			grpc_opentracing.UnaryServerInterceptor(tracing_opts...),
		),
	)
	pb.RegisterPersonServiceServer(grpcServer, &personService)
	pb.RegisterVehicleServiceServer(grpcServer, &vehicleService)

	ctx, cancel := context.WithCancel(ctx)
	go func() {
		defer cancel()
		logger.Infof("Grpc started!")
		if err = grpcServer.Serve(listener); err != nil {
			logger.Errorf(err.Error())
		}
	}()

	// GRPC http gateway
	gatewayServer := createGatewayServer(cfg.PersonService.Port, cfg.PersonService.GatewayPort)
	go func() {
		defer cancel()
		logger.Infof("Grpc gateway started!")
		if err := gatewayServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Errorf(err.Error())
		}
	}()

	// Так называемый, gracefull shutdown

	<-ctx.Done()
	log.Printf("ctx.Done\n")

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

	closer, err := c.InitGlobalTracer(cfg.PersonService.AppName)
	if err != nil {
		return nil, err
	}

	return closer, nil
}

type Closer func()

var NoCloser = func() {}

type Repos struct {
	person  repo.PersonRepo
	vehicle repo.VehicleRepo
}

func initRepo(cfg config.Config) (Repos, Closer) {
	if cfg.PersonService.Storage == config.StoragePostgres {
		pool, err := postgres.New(context.Background(), &cfg.Pooler)
		if err != nil {
			logger.FatalKV(err.Error())
		}

		return Repos{
			person:  postgres_repo.NewPersonRepo(pool),
			vehicle: postgres_repo.NewVehicleRepo(pool),
		}, NoCloser

	} else if cfg.PersonService.Storage == config.StorageMemory {
		return Repos{
			person:  memory_repo.NewPersonRepo(cfg.Storage),
			vehicle: memory_repo.NewVehicleRepo(cfg.Storage),
		}, NoCloser

	}

	logger.Warnf("Unsupported storage type %s. Using in-memory", cfg.PersonService.Storage)
	return Repos{
		person:  memory_repo.NewPersonRepo(cfg.Storage),
		vehicle: memory_repo.NewVehicleRepo(cfg.Storage),
	}, NoCloser
}
