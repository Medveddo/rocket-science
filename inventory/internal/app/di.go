package app

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	partApiV1 "github.com/Medveddo/rocket-science/inventory/internal/api/part/v1"
	"github.com/Medveddo/rocket-science/inventory/internal/config"
	"github.com/Medveddo/rocket-science/inventory/internal/repository"
	partRepository "github.com/Medveddo/rocket-science/inventory/internal/repository/part"
	"github.com/Medveddo/rocket-science/inventory/internal/service"
	partService "github.com/Medveddo/rocket-science/inventory/internal/service/part"
	"github.com/Medveddo/rocket-science/platform/pkg/closer"
	inventoryV1 "github.com/Medveddo/rocket-science/shared/pkg/proto/inventory/v1"
)

type diContainer struct {
	inventoryV1API inventoryV1.InventoryServiceServer

	partService    service.PartService
	partRepository repository.PartRepository

	mongoClient   *mongo.Client
	mongoDBHandle *mongo.Database
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) MongoDBClient(ctx context.Context) *mongo.Client {
	if d.mongoClient == nil {
		client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.AppConfig().Mongo.URI()))
		if err != nil {
			panic(fmt.Sprintf("failed to connect to MongoDB: %s\n", err.Error()))
		}

		err = client.Ping(ctx, nil)
		if err != nil {
			panic(fmt.Sprintf("failed to ping MongoDB: %v\n", err))
		}

		closer.AddNamed("MongoDB client", func(ctx context.Context) error {
			return client.Disconnect(ctx)
		})

		d.mongoClient = client
	}

	return d.mongoClient
}

func (d *diContainer) MongoDBHandle(ctx context.Context) *mongo.Database {
	if d.mongoDBHandle == nil {
		d.mongoDBHandle = d.MongoDBClient(ctx).Database(config.AppConfig().Mongo.DatabaseName())
	}

	return d.mongoDBHandle
}

func (d *diContainer) InventoryV1API(ctx context.Context) inventoryV1.InventoryServiceServer {
	if d.inventoryV1API == nil {
		d.inventoryV1API = partApiV1.NewPartAPI(d.PartService(ctx))
	}

	return d.inventoryV1API
}

func (d *diContainer) PartService(ctx context.Context) service.PartService {
	if d.partService == nil {
		d.partService = partService.NewPartService(d.PartRepository(ctx))
	}

	return d.partService
}

func (d *diContainer) PartRepository(ctx context.Context) repository.PartRepository {
	if d.partRepository == nil {
		repo := partRepository.NewPartRepository(d.MongoDBHandle(ctx))
		err := partRepository.InitParts(ctx, repo.Collection())
		if err != nil {
			panic(fmt.Sprintf("failed to init parts: %s\n", err.Error()))
		}
		d.partRepository = repo
	}

	return d.partRepository
}
