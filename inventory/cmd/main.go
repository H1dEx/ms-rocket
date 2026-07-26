package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	inventoryApi "github.com/H1dEx/ms-rocket/inventory/internal/api/inventory/v1"
	"github.com/H1dEx/ms-rocket/inventory/internal/config"
	inventoryRepo "github.com/H1dEx/ms-rocket/inventory/internal/repository/inventory"
	inventoryService "github.com/H1dEx/ms-rocket/inventory/internal/service/inventory"
	"github.com/H1dEx/ms-rocket/platform/pkg/grpc/health"
	inventoryV1 "github.com/H1dEx/ms-rocket/shared/pkg/proto/inventory/v1"
)

const configPath = "./deploy/compose/inventory/.env"

func main() {
	err := config.Load(configPath)
	if err != nil {
		log.Printf("failed to load config: %v\n", err)
		return
	}

	cfg := config.GetConfig()

	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.Mongo.URI()))
	if err != nil {
		log.Printf("failed to connect to MongoDB: %v\n", err)
		return
	}
	defer func() {
		if cerr := client.Disconnect(ctx); cerr != nil {
			log.Printf("failed to disconnect from MongoDB: %v\n", cerr)
		}
	}()

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Printf("failed to ping database: %v\n", err)
		return
	}
	conn := client.Database(cfg.Mongo.DatabaseName(), nil)

	lis, err := net.Listen("tcp", cfg.InventoryGRPC.Address())
	if err != nil {
		log.Printf("failed to listen: %v\n", err)
		return
	}

	defer func() {
		if cerr := lis.Close(); cerr != nil {
			log.Printf("failed to close listener: %v\n", cerr)
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
		log.Printf("🚀 gRPC server listening on %s\n", cfg.InventoryGRPC.Address())
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
