package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderV1 "github.com/Medveddo/rocket-science/shared/pkg/openapi/order/v1"
	paymentV1 "github.com/Medveddo/rocket-science/shared/pkg/proto/payment/v1"
)

const (
	httpPort = "8080"
	// Таймауты для HTTP-сервера
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second
)

const (
	ORDER_STATUS_UNKNOWN = iota
	ORDER_STATUS_PENDING_PAYMENT
	ORDER_STATUS_PAID
	ORDER_STATUS_CANCELED
)

type Order struct {
	OrderUuid       uuid.UUID
	UserUuid        uuid.UUID
	PartsUuids      []uuid.UUID
	TotalPrice      float64
	TransactionUuid uuid.UUID
	PaymentMethod   *string
	Status          uint8
}

type OrderStorage struct {
	mu     sync.RWMutex
	orders map[string]*Order
}

func NewOrderStorage() *OrderStorage {
	return &OrderStorage{
		orders: make(map[string]*Order),
	}
}

func (s *OrderStorage) UpdateOrder(order *Order) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.orders[order.OrderUuid.String()] = order

	return nil
}

type OrderHandler struct {
	storage       *OrderStorage
	paymentClient paymentV1.PaymentServiceClient
}

func NewOrderHandler(
	storage *OrderStorage,
	paymentClient paymentV1.PaymentServiceClient,
) *OrderHandler {
	return &OrderHandler{
		storage:       storage,
		paymentClient: paymentClient,
	}
}

func (h *OrderHandler) CreateOrder(ctx context.Context, req *orderV1.CreateOrderRequest) (orderV1.CreateOrderRes, error) {
	orderUuid := uuid.New()
	order := &Order{
		UserUuid:   req.UserUUID,
		OrderUuid:  orderUuid,
		PartsUuids: req.PartUuids,
		Status:     ORDER_STATUS_PENDING_PAYMENT,
	}
	err := h.storage.UpdateOrder(order)
	if err != nil {
		return nil, h.NewError(ctx, err)
	}

	response := &orderV1.CreateOrderResponse{
		OrderUUID:  order.OrderUuid,
		TotalPrice: 10,
	}
	return response, nil
}

func (h *OrderHandler) PayOrder(ctx context.Context, req *orderV1.PayOrderRequest, params orderV1.PayOrderParams) (orderV1.PayOrderRes, error) {
	h.storage.mu.Lock()
	defer h.storage.mu.Unlock()

	order, ok := h.storage.orders[params.OrderUUID]

	if !ok {
		return &orderV1.NotFoundError{
			Code:    404,
			Message: "Order with UUID: '" + params.OrderUUID + "' not found",
		}, nil
	}

	paymentMethodsMap := map[orderV1.PayOrderRequestPaymentMethod]paymentV1.PaymentMethod{
		orderV1.PayOrderRequestPaymentMethodCARD:          paymentV1.PaymentMethod_PAYMENT_METHOD_CARD,
		orderV1.PayOrderRequestPaymentMethodSBP:           paymentV1.PaymentMethod_PAYMENT_METHOD_SBP,
		orderV1.PayOrderRequestPaymentMethodCREDITCARD:    paymentV1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD,
		orderV1.PayOrderRequestPaymentMethodINVESTORMONEY: paymentV1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY,
	}

	paymentMethod, ok := paymentMethodsMap[req.PaymentMethod]
	if !ok {
		return &orderV1.BadRequestError{
			Code:    400,
			Message: "Payment method '" + string(req.PaymentMethod) + "' is not supported",
		}, nil
	}

	response, err := h.paymentClient.PayOrder(ctx, &paymentV1.PayOrderRequest{
		UserUuid:      order.UserUuid.String(),
		OrderUuid:     order.OrderUuid.String(),
		PaymentMethod: paymentMethod,
	})
	if err != nil {
		return &orderV1.InternalServerError{
			Code:    500,
			Message: "failed to process payment",
		}, err
	}

	transactionUUID, err := uuid.Parse(response.TransactionUuid)
	if err != nil {
		return &orderV1.InternalServerError{
			Code:    500,
			Message: "failed to process payment",
		}, err
	}

	paymentMethodStr := string(req.PaymentMethod)
	order.Status = ORDER_STATUS_PAID
	order.PaymentMethod = &paymentMethodStr
	order.TransactionUuid = transactionUUID

	return &orderV1.PayOrderResponse{
		TransactionUUID: transactionUUID,
	}, nil
}

// NewError создает новую ошибку в формате GenericError
func (h *OrderHandler) NewError(_ context.Context, err error) *orderV1.GenericErrorStatusCode {
	return &orderV1.GenericErrorStatusCode{
		StatusCode: http.StatusInternalServerError,
		Response: orderV1.GenericError{
			Code:    orderV1.NewOptInt(http.StatusInternalServerError),
			Message: orderV1.NewOptString(err.Error()),
		},
	}
}

func main() {
	// Создаем хранилище для данных о погоде
	storage := NewOrderStorage()

	conn, err := grpc.NewClient(
		"localhost:50052",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		log.Fatalf("failed to connect to Payment Service: %v\n", err)
	}

	// to warm up channels
	conn.Connect()

	paymentClient := paymentV1.NewPaymentServiceClient(conn)

	// Создаем обработчик API погоды
	orderHandler := NewOrderHandler(storage, paymentClient)

	// Создаем OpenAPI сервер
	orderServer, err := orderV1.NewServer(orderHandler)
	if err != nil {
		log.Fatalf("ошибка создания сервера OpenAPI: %v", err)
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
