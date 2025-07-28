package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	partApiV1 "github.com/Medveddo/rocket-science/inventory/internal/api/part/v1"
	"github.com/Medveddo/rocket-science/inventory/internal/config"
	partRepository "github.com/Medveddo/rocket-science/inventory/internal/repository/part"
	partService "github.com/Medveddo/rocket-science/inventory/internal/service/part"
	"github.com/Medveddo/rocket-science/shared/pkg/interceptor"
	inventoryV1 "github.com/Medveddo/rocket-science/shared/pkg/proto/inventory/v1"
)

const configPath = "../deploy/compose/inventory/.env"

func main() {
	ctx := context.Background()

	err := config.Load(configPath)
	if err != nil {
		log.Printf("cannot load config: %v\n", err)
		return
	}

	mongoURI := config.AppConfig().Mongo.URI()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Printf("Ошибка подключения к MongoDB: %v\n", err)
		return
	}

	defer func() {
		if cerr := client.Disconnect(ctx); cerr != nil {
			log.Printf("Ошибка при отключении от MongoDB: %v\n", cerr)
		}
	}()

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Printf("MongoDB недоступна, ошибка ping: %v\n", err)
		return
	}
	log.Println("Успешное подключение к MongoDB")

	inventoryMongoDb := client.Database("inventory")

	repository, err := partRepository.NewPartRepository(inventoryMongoDb)
	if err != nil {
		log.Printf("error while initializing part repository: %v\n", err)
		return
	}
	service := partService.NewPartService(repository)
	api := partApiV1.NewPartAPI(service)

	grpcAddress := config.AppConfig().InventoryGRPC.Address()
	lis, err := net.Listen("tcp", grpcAddress)
	if err != nil {
		log.Printf("failed to listen: %v\n", err)
		return
	}
	defer func() {
		if cerr := lis.Close(); cerr != nil {
			log.Printf("failed to close listener: %v\n", cerr)
		}
	}()

	// Создаем gRPC сервер
	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpc.UnaryServerInterceptor(interceptor.LoggerInterceptor()),
			recovery.UnaryServerInterceptor(),
		),
	)

	inventoryV1.RegisterInventoryServiceServer(s, api)

	// Включаем рефлексию для отладки
	reflection.Register(s)

	go func() {
		log.Printf("🚀 gRPC server listening on %s\n", grpcAddress)
		err = s.Serve(lis)
		if err != nil {
			log.Printf("failed to serve: %v\n", err)
			return
		}
	}()

	// Запускаем HTTP сервер с gRPC Gateway и Swagger UI
	var gwServer *http.Server
	go func() {
		// Создаем контекст с отменой
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// Создаем мультиплексор для HTTP запросов
		mux := runtime.NewServeMux()

		// Настраиваем опции для соединения с gRPC сервером
		opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

		// Регистрируем gRPC-gateway хендлеры
		err = inventoryV1.RegisterInventoryServiceHandlerFromEndpoint(
			ctx,
			mux,
			grpcAddress,
			opts,
		)
		if err != nil {
			log.Printf("Failed to register gateway: %v\n", err)
			return
		}

		// Создаем файловый сервер для swagger-ui
		fileServer := http.FileServer(http.Dir("api"))

		// Создаем HTTP маршрутизатор
		httpMux := http.NewServeMux()

		// Регистрируем API эндпоинты
		httpMux.Handle("/api/", mux)

		// Swagger UI эндпоинты
		httpMux.Handle("/swagger-ui.html", fileServer)
		httpMux.Handle("/swagger.json", fileServer)

		// Редирект с корня на Swagger UI
		httpMux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/" {
				http.Redirect(w, r, "/swagger-ui.html", http.StatusMovedPermanently)
				return
			}
			fileServer.ServeHTTP(w, r)
		}))

		// Создаем HTTP сервер
		httpAddr := config.AppConfig().HTTP.Address()
		httpPort := config.AppConfig().HTTP.Port()

		gwServer = &http.Server{
			Addr:              httpAddr,
			Handler:           httpMux,
			ReadHeaderTimeout: 10 * time.Second,
		}

		// Запускаем HTTP сервер
		log.Printf("🌐 HTTP server with gRPC-Gateway and Swagger UI listening on %s\n", httpPort)
		err = gwServer.ListenAndServe()
		if err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Printf("Failed to serve HTTP: %v\n", err)
			return
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down gRPC server...")
	s.GracefulStop()
	log.Println("✅ Server stopped")
}
