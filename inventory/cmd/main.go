package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	// "github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	"github.com/Medveddo/rocket-science/shared/pkg/interceptor"
	inventoryV1 "github.com/Medveddo/rocket-science/shared/pkg/proto/inventory/v1"
)

const grpcPort = 50051

type inventoryService struct {
	inventoryV1.UnimplementedInventoryServiceServer

	mu    sync.RWMutex
	parts map[string]*inventoryV1.Part
}

func NewInventoryService() *inventoryService {
	service := &inventoryService{
		parts: make(map[string]*inventoryV1.Part),
	}
	part1 := inventoryV1.Part{
		Uuid:          "111e4567-e89b-12d3-a456-426614174001",
		Name:          "Hyperdrive Engine",
		Description:   "A class-9 hyperdrive engine capable of faster-than-light travel.",
		Price:         450000.00,
		StockQuantity: 3,
		Category:      inventoryV1.Category_CATEGORY_ENGINE,
	}

	part2 := inventoryV1.Part{
		Uuid:          "222e4567-e89b-12d3-a456-426614174002",
		Name:          "Quantum Shield Generator",
		Description:   "Advanced shield generator providing protection against cosmic radiation.",
		Price:         175000.00,
		StockQuantity: 5,
		Category:      inventoryV1.Category_CATEGORY_SHIELD,
	}

	service.parts[part1.Uuid] = &part1
	service.parts[part2.Uuid] = &part2
	return service
}

func (s *inventoryService) GetPart(_ context.Context, request *inventoryV1.GetPartRequest) (*inventoryV1.GetPartResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	part, ok := s.parts[request.GetUuid()]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "sighting with UUID %s not found", request.GetUuid())
	}

	return &inventoryV1.GetPartResponse{
		Part: part,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
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
		),
	)

	// Регистрируем наш сервис
	service := NewInventoryService()

	inventoryV1.RegisterInventoryServiceServer(s, service)

	// Включаем рефлексию для отладки
	reflection.Register(s)

	go func() {
		log.Printf("🚀 gRPC server listening on %d\n", grpcPort)
		err = s.Serve(lis)
		if err != nil {
			log.Printf("failed to serve: %v\n", err)
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
