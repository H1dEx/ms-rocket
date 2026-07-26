package main

import (
	"context"
	"fmt"
	"net"
	"os/signal"
	"syscall"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	inventoryApi "github.com/H1dEx/ms-rocket/inventory/internal/api/inventory/v1"
	"github.com/H1dEx/ms-rocket/inventory/internal/config"
	inventoryRepo "github.com/H1dEx/ms-rocket/inventory/internal/repository/inventory"
	inventoryService "github.com/H1dEx/ms-rocket/inventory/internal/service/inventory"
	"github.com/H1dEx/ms-rocket/platform/pkg/grpc/health"
	"github.com/H1dEx/ms-rocket/platform/pkg/logger"
	inventoryV1 "github.com/H1dEx/ms-rocket/shared/pkg/proto/inventory/v1"
)

const configPath = "./deploy/compose/inventory/.env"

func main() {
	appCtx, appCancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer appCancel()

	err := logger.Init("info", true)
	if err != nil {
		fmt.Println("failed to init logger", err)
		return
	}
	err = config.Load(configPath)
	if err != nil {
		logger.Error(appCtx, "failed to load config", zap.Error(err))
		return
	}

	cfg := config.GetConfig()

	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.Mongo.URI()))
	if err != nil {
		logger.Error(appCtx, "failed to connect to MongoDB", zap.Error(err))
		return
	}
	defer func() {
		if cerr := client.Disconnect(ctx); cerr != nil {
			logger.Error(appCtx, "failed to disconnect from MongoDB", zap.Error(cerr))
		}
	}()

	err = client.Ping(ctx, nil)
	if err != nil {
		logger.Error(appCtx, "failed to ping database", zap.Error(err))
		return
	}
	conn := client.Database(cfg.Mongo.DatabaseName(), nil)

	lis, err := net.Listen("tcp", cfg.InventoryGRPC.Address())
	if err != nil {
		logger.Error(appCtx, "failed to listen", zap.Error(err))
		return
	}

	defer func() {
		if cerr := lis.Close(); cerr != nil {
			logger.Error(appCtx, "failed to close listener", zap.Error(cerr))
		}
	}()

	repo := inventoryRepo.NewRepository(conn)
	service := inventoryService.NewService(repo)
	api := inventoryApi.NewAPI(service)

	s := grpc.NewServer()

	health.RegisterServer(s)

	inventoryV1.RegisterInventoryServiceServer(s, api)
	reflection.Register(s)

	go func() {
		logger.Info(appCtx, fmt.Sprintf("🚀 gRPC server listening on %s", cfg.InventoryGRPC.Address()))
		err = s.Serve(lis)
		if err != nil {
			logger.Error(appCtx, "failed to serve", zap.Error(err))
			return
		}
	}()

	// Graceful shutdown
	<-appCtx.Done()
	logger.Info(appCtx, "🛑 Shutting down gRPC server...")
	s.GracefulStop()
	logger.Info(appCtx, "✅ Server stopped")
}
