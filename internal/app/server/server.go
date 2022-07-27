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

	"gitlab.ozon.dev/igor.benko.1991/homework/internal/config"
	"gitlab.ozon.dev/igor.benko.1991/homework/internal/service"
	"gitlab.ozon.dev/igor.benko.1991/homework/internal/storage"
	pb "gitlab.ozon.dev/igor.benko.1991/homework/pkg/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

const _shotdownTimeout = 5 * time.Second

func Run() {
	cfg := config.Get()
	ctx, cancel := context.WithCancel(context.Background())

	// Инициализация хранилища
	memoryStorage := storage.NewMemoryStorage()

	// Инициализация сервиса
	personService := service.NewPersonService(memoryStorage)

	// GRPC
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Grpc.Port))
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterPersonServiceServer(grpcServer, &personService)

	go func() {
		log.Println("Grpc started!")
		if err = grpcServer.Serve(listener); err != nil {
			log.Println(err)
			cancel()
		}
	}()

	// GRPC http gateway
	gatewayServer := createGatewayServer(cfg.Grpc.Port, cfg.Grpc.GatewayPort)
	go func() {
		log.Println("Grpc gateway started!")
		if err := gatewayServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Println(err)
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
	ctx, cancel = context.WithTimeout(context.Background(), _shotdownTimeout)
	defer cancel()

	if err := gatewayServer.Shutdown(ctx); err != nil {
		log.Printf("Error on gatewayServer stop: %s\n", err)
	} else {
		log.Println("gatewayServer stopped")
	}

	grpcServer.GracefulStop()
	log.Println("Grpc server stopped")
}

func createGatewayServer(grpcPort, httpPort int) *http.Server {
	ctx := context.Background()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	if err := pb.RegisterPersonServiceHandlerFromEndpoint(ctx, mux, fmt.Sprintf(":%d", grpcPort), opts); err != nil {
		log.Fatal(err)
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Group("v1/*{grpc_gateway}").Any("", gin.WrapH(mux))
	router.StaticFile("/docs/person.swagger.json", "./docs/person.swagger.json")
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/docs/person.swagger.json")))

	return &http.Server{
		Addr:    fmt.Sprintf(":%d", httpPort),
		Handler: router,
	}
}
