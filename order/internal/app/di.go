package app

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderApiV1 "github.com/Medveddo/rocket-science/order/internal/api/order/v1"
	grpcClient "github.com/Medveddo/rocket-science/order/internal/client/grpc"
	inventoryClientV1 "github.com/Medveddo/rocket-science/order/internal/client/grpc/inventory/v1"
	paymentClientV1 "github.com/Medveddo/rocket-science/order/internal/client/grpc/payment/v1"
	"github.com/Medveddo/rocket-science/order/internal/config"
	"github.com/Medveddo/rocket-science/order/internal/repository"
	orderRepository "github.com/Medveddo/rocket-science/order/internal/repository/order"
	"github.com/Medveddo/rocket-science/order/internal/service"
	orderService "github.com/Medveddo/rocket-science/order/internal/service/order"
	"github.com/Medveddo/rocket-science/platform/pkg/closer"
	"github.com/Medveddo/rocket-science/platform/pkg/migrator"
	pgMigrator "github.com/Medveddo/rocket-science/platform/pkg/migrator/pg"
	orderV1 "github.com/Medveddo/rocket-science/shared/pkg/openapi/order/v1"
	inventoryV1 "github.com/Medveddo/rocket-science/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/Medveddo/rocket-science/shared/pkg/proto/payment/v1"
)

type diContainer struct {
	orderV1API orderV1.Handler

	orderService    service.OrderService
	orderRepository repository.OrderRepository

	paymentClient   grpcClient.PaymentClient
	inventoryClient grpcClient.InventoryClient

	paymentConn   *grpc.ClientConn
	inventoryConn *grpc.ClientConn
	postgresPool  *pgxpool.Pool

	migrator migrator.Migrator
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) PostgreSQLPool(ctx context.Context) *pgxpool.Pool {
	if d.postgresPool == nil {
		pool, err := pgxpool.New(ctx, config.AppConfig().PostgreSQL.URI())
		if err != nil {
			panic(fmt.Sprintf("failed to connect to PostgreSQL: %s\n", err.Error()))
		}

		err = pool.Ping(ctx)
		if err != nil {
			panic(fmt.Sprintf("failed to ping PostgreSQL: %v\n", err))
		}

		closer.AddNamed("PostgreSQL pool", func(ctx context.Context) error {
			pool.Close()
			return nil
		})

		d.postgresPool = pool
	}

	return d.postgresPool
}

func (d *diContainer) PaymentConnection(ctx context.Context) *grpc.ClientConn {
	if d.paymentConn == nil {
		conn, err := grpc.NewClient(
			config.AppConfig().PaymentGRPC.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to connect to Payment Service: %s\n", err.Error()))
		}

		closer.AddNamed("Payment gRPC connection", func(ctx context.Context) error {
			return conn.Close()
		})

		d.paymentConn = conn
	}

	return d.paymentConn
}

func (d *diContainer) InventoryConnection(ctx context.Context) *grpc.ClientConn {
	if d.inventoryConn == nil {
		conn, err := grpc.NewClient(
			config.AppConfig().InventoryGRPC.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to connect to Inventory Service: %s\n", err.Error()))
		}

		closer.AddNamed("Inventory gRPC connection", func(ctx context.Context) error {
			return conn.Close()
		})

		d.inventoryConn = conn
	}

	return d.inventoryConn
}

func (d *diContainer) OrderV1API(ctx context.Context) orderV1.Handler {
	if d.orderV1API == nil {
		d.orderV1API = orderApiV1.NewOrderAPI(d.OrderService(ctx))
	}

	return d.orderV1API
}

func (d *diContainer) OrderService(ctx context.Context) service.OrderService {
	if d.orderService == nil {
		d.orderService = orderService.NewOrderService(
			d.OrderRepository(ctx),
			d.InventoryClient(ctx),
			d.PaymentClient(ctx),
		)
	}

	return d.orderService
}

func (d *diContainer) OrderRepository(ctx context.Context) repository.OrderRepository {
	if d.orderRepository == nil {
		d.orderRepository = orderRepository.NewOrderRepository(d.PostgreSQLPool(ctx))
	}

	return d.orderRepository
}

func (d *diContainer) PaymentClient(ctx context.Context) grpcClient.PaymentClient {
	if d.paymentClient == nil {
		paymentProtoClient := paymentV1.NewPaymentServiceClient(d.PaymentConnection(ctx))
		d.paymentClient = paymentClientV1.NewPaymentClientV1(paymentProtoClient)
	}

	return d.paymentClient
}

func (d *diContainer) InventoryClient(ctx context.Context) grpcClient.InventoryClient {
	if d.inventoryClient == nil {
		inventoryProtoClient := inventoryV1.NewInventoryServiceClient(d.InventoryConnection(ctx))
		d.inventoryClient = inventoryClientV1.NewInventoryClientV1(inventoryProtoClient)
	}

	return d.inventoryClient
}

func (d *diContainer) Migrator(ctx context.Context) migrator.Migrator {
	if d.migrator == nil {
		migrationsDir := "./migrations"
		db := stdlib.OpenDB(*d.PostgreSQLPool(ctx).Config().ConnConfig)
		d.migrator = pgMigrator.NewMigrator(db, migrationsDir)
	}

	return d.migrator
}
