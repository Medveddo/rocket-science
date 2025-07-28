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

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderApiV1 "github.com/Medveddo/rocket-science/order/internal/api/order/v1"
	inventoryClientV1 "github.com/Medveddo/rocket-science/order/internal/client/grpc/inventory/v1"
	paymentClientV1 "github.com/Medveddo/rocket-science/order/internal/client/grpc/payment/v1"
	orderRepository "github.com/Medveddo/rocket-science/order/internal/repository/order"
	orderService "github.com/Medveddo/rocket-science/order/internal/service/order"
	orderV1 "github.com/Medveddo/rocket-science/shared/pkg/openapi/order/v1"
	inventoryV1 "github.com/Medveddo/rocket-science/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/Medveddo/rocket-science/shared/pkg/proto/payment/v1"
)

const (
	httpPort = "8080"
	// Таймауты для HTTP-сервера
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second
)

func main() {
	ctx := context.Background()
	err := godotenv.Load("../.env")
	if err != nil {
		log.Printf("cannot load environment file: %v\n", err)
		return
	}

	dbURI := os.Getenv("POSTGRES_URI")
	if dbURI == "" {
		log.Printf("error: got empty MONGO_URI variable \n")
		return
	}

	pool, err := pgxpool.New(ctx, dbURI)
	if err != nil {
		log.Printf("failed to connect to database: %v\n", err)
		return
	}
	defer pool.Close()
	err = pool.Ping(ctx)
	if err != nil {
		log.Printf("failed to ping database: %v\n", err)
		return
	}

	conn, err := grpc.NewClient(
		"localhost:50052",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("failed to connect to Payment Service: %v\n", err)
		return
	}

	inventoryConn, err := grpc.NewClient(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("failed to connect to Payment Service: %v\n", err)
		return
	}

	// to warm up channels
	conn.Connect()

	paymentProtoClient := paymentV1.NewPaymentServiceClient(conn)
	inventoryProtoClient := inventoryV1.NewInventoryServiceClient(inventoryConn)

	paymentClient := paymentClientV1.NewPaymentClientV1(paymentProtoClient)
	inventoryClient := inventoryClientV1.NewInventoryClientV1(inventoryProtoClient)

	repository := orderRepository.NewOrderRepository(pool)
	service := orderService.NewOrderService(repository, inventoryClient, paymentClient)
	api := orderApiV1.NewOrderAPI(service)

	// Создаем OpenAPI сервер
	orderServer, err := orderV1.NewServer(api)
	if err != nil {
		log.Printf("ошибка создания сервера OpenAPI: %v", err)
		return
	}

	// Инициализируем роутер Chi
	r := chi.NewRouter()

	// Добавляем middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))
	// r.Use(customMiddleware.RequestLogger)

	// Монтируем обработчики OpenAPI
	r.Mount("/", orderServer)

	// Запускаем HTTP-сервер
	server := &http.Server{
		Addr:              net.JoinHostPort("localhost", httpPort),
		Handler:           r,
		ReadHeaderTimeout: readHeaderTimeout, // Защита от Slowloris атак - тип DDoS-атаки, при которой
		// атакующий умышленно медленно отправляет HTTP-заголовки, удерживая соединения открытыми и истощая
		// пул доступных соединений на сервере. ReadHeaderTimeout принудительно закрывает соединение,
		// если клиент не успел отправить все заголовки за отведенное время.
	}

	// Запускаем сервер в отдельной горутине
	go func() {
		log.Printf("🚀 HTTP-сервер запущен на порту %s\n", httpPort)
		err = server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("❌ Ошибка запуска сервера: %v\n", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Завершение работы сервера...")

	// Создаем контекст с таймаутом для остановки сервера
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Printf("❌ Ошибка при остановке сервера: %v\n", err)
	}

	log.Println("✅ Сервер остановлен")
}
