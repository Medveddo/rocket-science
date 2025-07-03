package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"slices"
	"sync"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

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

	now := time.Now()

	part1 := inventoryV1.Part{
		Uuid:          "111e4567-e89b-12d3-a456-426614174001",
		Name:          "Hyperdrive Engine",
		Description:   "A class-9 hyperdrive engine capable of faster-than-light travel.",
		Price:         450000.00,
		StockQuantity: 3,
		Category:      inventoryV1.Category_CATEGORY_ENGINE,
		Dimensions: &inventoryV1.Dimensions{
			Length: 120.0,
			Width:  80.0,
			Height: 100.0,
			Weight: 500.0,
		},
		Manufacturer: &inventoryV1.Manufacturer{
			Name:    "Hyperdrive Corp",
			Country: "USA",
			Website: "https://hyperdrive.example.com",
		},
		Tags: []string{"engine", "hyperdrive", "space"},
		Metadata: map[string]*inventoryV1.Value{
			"power_output":    {Kind: &inventoryV1.Value_DoubleValue{DoubleValue: 9.5}},
			"is_experimental": {Kind: &inventoryV1.Value_BoolValue{BoolValue: true}},
		},
		CreatedAt: timestamppb.New(now),
		UpdatedAt: timestamppb.New(now),
	}

	part2 := inventoryV1.Part{
		Uuid:          "222e4567-e89b-12d3-a456-426614174002",
		Name:          "Quantum Shield Generator",
		Description:   "Advanced shield generator providing protection against cosmic radiation.",
		Price:         175000.00,
		StockQuantity: 5,
		Category:      inventoryV1.Category_CATEGORY_SHIELD,
		Dimensions: &inventoryV1.Dimensions{
			Length: 60.0,
			Width:  40.0,
			Height: 50.0,
			Weight: 150.0,
		},
		Manufacturer: &inventoryV1.Manufacturer{
			Name:    "Quantum Tech",
			Country: "Germany",
			Website: "https://quantumtech.example.com",
		},
		Tags: []string{"shield", "quantum", "defense"},
		Metadata: map[string]*inventoryV1.Value{
			"energy_consumption": {Kind: &inventoryV1.Value_DoubleValue{DoubleValue: 3.2}},
			"warranty_years":     {Kind: &inventoryV1.Value_Int64Value{Int64Value: 5}},
		},
		CreatedAt: timestamppb.New(now),
		UpdatedAt: timestamppb.New(now),
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

func (s *inventoryService) ListParts(_ context.Context, request *inventoryV1.ListPartsRequest) (*inventoryV1.ListPartsResponse, error) {
	parts := []*inventoryV1.Part{}

	for _, v := range s.parts {
		// Если передан фильтр по UUID
		if len(request.Filter.Uuids) > 0 {
			// Если UUID детали нет в фильтре - continue
			if !slices.Contains(request.Filter.Uuids, v.Uuid) {
				continue
			}
		}
		if len(request.Filter.Names) > 0 {
			if !slices.Contains(request.Filter.Names, v.Name) {
				continue
			}
		}
		if len(request.Filter.Categories) > 0 {
			if !slices.Contains(request.Filter.Categories, v.Category) {
				continue
			}
		}
		if len(request.Filter.ManufacturerCountries) > 0 {
			if !slices.Contains(request.Filter.ManufacturerCountries, v.Manufacturer.Country) {
				continue
			}
		}
		if len(request.Filter.Tags) > 0 {
			needToContinue := false
			for _, tag := range request.Filter.Tags {
				if !slices.Contains(v.Tags, tag) {
					needToContinue = true
					break
				}
			}
			if needToContinue {
				continue
			}

		}
		parts = append(parts, v)
	}
	response := inventoryV1.ListPartsResponse{
		Parts: parts,
	}
	return &response, nil
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
