package app

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"

	inventoryApi "github.com/H1dEx/ms-rocket/inventory/internal/api/inventory/v1"
	"github.com/H1dEx/ms-rocket/inventory/internal/config"
	"github.com/H1dEx/ms-rocket/inventory/internal/repository"
	inventoryRepo "github.com/H1dEx/ms-rocket/inventory/internal/repository/inventory"
	"github.com/H1dEx/ms-rocket/inventory/internal/service"
	inventoryService "github.com/H1dEx/ms-rocket/inventory/internal/service/inventory"
	"github.com/H1dEx/ms-rocket/platform/pkg/closer"
	"github.com/H1dEx/ms-rocket/platform/pkg/logger"
	inventoryV1 "github.com/H1dEx/ms-rocket/shared/pkg/proto/inventory/v1"
)

type diContainer struct {
	inventoryV1API   inventoryV1.InventoryServiceServer
	inventoryService service.InventoryService
	inventoryRepo    repository.InventoryRepository

	mongoDBClient *mongo.Client
	mongoDBHandle *mongo.Database
}

func NewDIContainer() *diContainer {
	return &diContainer{}
}

func (c *diContainer) InventoryV1API(ctx context.Context) inventoryV1.InventoryServiceServer {
	if c.inventoryV1API == nil {
		c.inventoryV1API = inventoryApi.NewAPI(c.InventoryService(ctx))
	}
	return c.inventoryV1API
}

func (c *diContainer) InventoryService(ctx context.Context) service.InventoryService {
	if c.inventoryService == nil {
		c.inventoryService = inventoryService.NewService(c.InventoryRepo(ctx))
	}
	return c.inventoryService
}

func (c *diContainer) InventoryRepo(ctx context.Context) repository.InventoryRepository {
	if c.inventoryRepo == nil {
		//nolint:contextcheck // We no need to pass context to the repository
		c.inventoryRepo = inventoryRepo.NewRepository(c.MongoDBHandle(ctx))
	}
	return c.inventoryRepo
}

func (c *diContainer) MongoDBHandle(ctx context.Context) *mongo.Database {
	if c.mongoDBHandle == nil {
		d := c.MongoDBClient(ctx).Database(config.GetConfig().Mongo.DatabaseName())
		c.mongoDBHandle = d
	}
	return c.mongoDBHandle
}

func (c *diContainer) MongoDBClient(ctx context.Context) *mongo.Client {
	if c.mongoDBClient == nil {
		client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.GetConfig().Mongo.URI()))
		if err != nil {
			panic(fmt.Errorf("failed to connect to MongoDB: %s", err.Error()))
		}
		closer.AddNamed("MongoDB connection", func(ctx context.Context) error {
			if cerr := client.Disconnect(ctx); cerr != nil {
				logger.Error(ctx, "failed to disconnect from MongoDB", zap.Error(cerr))
				return cerr
			}
			return nil
		})

		err = client.Ping(ctx, nil)
		if err != nil {
			panic(fmt.Errorf("failed to ping database: %s", err.Error()))
		}
		c.mongoDBClient = client
	}
	return c.mongoDBClient
}
