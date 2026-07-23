package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	inventoryApi "github.com/H1dEx/ms-rocket/inventory/internal/api/inventory/v1"
	inventoryRepo "github.com/H1dEx/ms-rocket/inventory/internal/repository/inventory"
	inventoryService "github.com/H1dEx/ms-rocket/inventory/internal/service/inventory"
	inventoryV1 "github.com/H1dEx/ms-rocket/shared/pkg/proto/inventory/v1"
)

const grpcPort = 50052

func main() {
	envFile, err := findEnvFile()
	if err != nil {
		log.Printf("failed to find .env file: %v\n", err)
		return
	}

	err = godotenv.Load(envFile)
	if err != nil {
		log.Printf("failed to load .env file: %v\n", err)
		return
	}

	INVENTORY_MONGO_URI := os.Getenv("INVENTORY_MONGO_URI")
	if INVENTORY_MONGO_URI == "" {
		log.Printf("INVENTORY_MONGO_URI is not set")
		return
	}

	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(INVENTORY_MONGO_URI))
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

	_ = client.Database("inventory-service").Collection("parts")

	// test := m.PartMongo{
	// 	UUID:          "123",
	// 	Name:          "Test",
	// 	Description:   "Test",
	// 	Price:         100,
	// 	StockQuantity: 100,
	// 	Category:      "Test",
	// 	Dimensions:    m.Dimensions{Width: 100, Height: 100},
	// }

	// _, err = db.InsertOne(ctx, test)
	// if err != nil {
	// 	log.Printf("failed to insert test: %v\n", err)
	// 	return
	// }

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

	repo := inventoryRepo.NewRepository()
	service := inventoryService.NewService(repo)
	api := inventoryApi.NewApi(service)

	s := grpc.NewServer()

	inventoryV1.RegisterInventoryServiceServer(s, api)
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

func findEnvFile() (string, error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", errors.New("runtime.Caller failed")
	}
	dir := filepath.Dir(file)
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.work")); err == nil {
			break
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", errors.New("go.work not found")
		}
		dir = parent
	}
	return filepath.Join(dir, ".env"), nil
}
