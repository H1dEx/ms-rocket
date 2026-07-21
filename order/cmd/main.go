package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"
	"time"

	"github.com/H1dEx/ms-rocket/order/internal/migrator"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderApi "github.com/H1dEx/ms-rocket/order/internal/api/order/v1"
	inventoryCli "github.com/H1dEx/ms-rocket/order/internal/client/grpc/inventory/v1"
	paymentCli "github.com/H1dEx/ms-rocket/order/internal/client/grpc/payment/v1"
	orderRepo "github.com/H1dEx/ms-rocket/order/internal/repository/order"
	orderService "github.com/H1dEx/ms-rocket/order/internal/service/order"
	orderV1 "github.com/H1dEx/ms-rocket/shared/pkg/openapi/order/v1"
	inventoryV1 "github.com/H1dEx/ms-rocket/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/H1dEx/ms-rocket/shared/pkg/proto/payment/v1"
)

const (
	httpPort = "8080"
	// Таймауты для HTTP-сервера
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second
	paymentPort       = "localhost:50051"
	inventoryPort     = "localhost:50052"
)

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

	DB_URI := os.Getenv("DB_URI")
	if DB_URI == "" {
		log.Printf("DB_URI is not set")
		return
	}

	orderDir, err := findOrderDir()
	if err != nil {
		log.Printf("failed to find order directory: %v\n", err)
		return
	}

	MIGRATION_PATH := os.Getenv("MIGRATIONS_DIR")
	if MIGRATION_PATH == "" {
		log.Printf("MIGRATIONS_DIR is not set")
		return
	}

	migrationDir := filepath.Join(orderDir, MIGRATION_PATH)
	ctx := context.Background()
	conn, err := pgxpool.New(ctx, DB_URI)
	if err != nil {
		log.Printf("failed to connect to database: %v\n", err)
		return
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err = conn.Ping(ctx)
	if err != nil {
		log.Printf("failed to ping database: %v\n", err)
		return
	}

	migrator := migrator.NewMigrator(stdlib.OpenDB(*conn.Config().ConnConfig), migrationDir)
	err = migrator.Up()
	if err != nil {
		log.Printf("failed to migrate database: %v\n", err)
		return
	}

	paymentConn, err := grpc.NewClient(
		paymentPort,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("failed to connect: %v\n", err)
		return
	}
	defer func() {
		if cerr := paymentConn.Close(); cerr != nil {
			log.Printf("failed to close connect: %v", cerr)
		}
	}()

	paymentClient := paymentV1.NewPaymentServiceClient(paymentConn)

	inventoryConn, err := grpc.NewClient(
		inventoryPort,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("failed to connect: %v\n", err)
		return
	}
	defer func() {
		if cerr := inventoryConn.Close(); cerr != nil {
			log.Printf("failed to close connect: %v", cerr)
		}
	}()

	inventoryClient := inventoryV1.NewInventoryServiceClient(inventoryConn)
	repo := orderRepo.NewOrderRepository()
	service := orderService.NewOrderService(repo, inventoryCli.NewInventoryClient(inventoryClient), paymentCli.MewPaymentClient(paymentClient))
	api := orderApi.NewOrderApi(service)

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

	r.Mount("/", orderServer)

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
	ctx, cancel = context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Printf("❌ Ошибка при остановке сервера: %v\n", err)
	}

	log.Println("✅ Сервер остановлен")
}

func findOrderDir() (string, error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", errors.New("runtime.Caller failed")
	}
	dir := filepath.Dir(file)
	for {
		if filepath.Base(dir) == "order" {
			break
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", errors.New("order not found")
		}
		dir = parent
	}
	return dir, nil
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
