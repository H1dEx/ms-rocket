package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderApi "github.com/H1dEx/ms-rocket/order/internal/api/order/v1"
	inventoryCli "github.com/H1dEx/ms-rocket/order/internal/client/grpc/inventory/v1"
	paymentCli "github.com/H1dEx/ms-rocket/order/internal/client/grpc/payment/v1"
	"github.com/H1dEx/ms-rocket/order/internal/config"
	"github.com/H1dEx/ms-rocket/order/internal/migrator"
	orderRepo "github.com/H1dEx/ms-rocket/order/internal/repository/order"
	orderService "github.com/H1dEx/ms-rocket/order/internal/service/order"
	orderV1 "github.com/H1dEx/ms-rocket/shared/pkg/openapi/order/v1"
	inventoryV1 "github.com/H1dEx/ms-rocket/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/H1dEx/ms-rocket/shared/pkg/proto/payment/v1"
)

const (
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second
	configPath        = "./deploy/compose/order/.env"
)

func main() {
	err := config.Load(configPath)
	if err != nil {
		log.Printf("failed to load config: %v\n", err)
		return
	}

	orderDir, err := findOrderDir()
	if err != nil {
		log.Printf("failed to find order directory: %v\n", err)
		return
	}

	ctx := context.Background()
	conn, err := pgxpool.New(ctx, config.GetConfig().Postgres.URI())
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
	sqlDB := stdlib.OpenDB(*conn.Config().ConnConfig)
	defer func() {
		if cerr := sqlDB.Close(); cerr != nil {
			log.Printf("failed to close database: %v", cerr)
		}
	}()
	migrationDir := filepath.Join(orderDir, config.GetConfig().Postgres.MigrationDir())
	migrator := migrator.NewMigrator(sqlDB, migrationDir)
	err = migrator.Up()
	if err != nil {
		log.Printf("failed to migrate database: %v\n", err)
		return
	}

	paymentConn, err := grpc.NewClient(
		config.GetConfig().PaymentGRPC.Address(),
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
		config.GetConfig().InventoryGRPC.Address(),
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
	repo := orderRepo.NewOrderRepository(conn)
	service := orderService.NewOrderService(repo, inventoryCli.NewInventoryClient(inventoryClient), paymentCli.NewPaymentClient(paymentClient))
	api := orderApi.NewOrderAPI(service)

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
		Addr:              config.GetConfig().OrderHTTP.Address(),
		Handler:           r,
		ReadHeaderTimeout: readHeaderTimeout, // Защита от Slowloris атак - тип DDoS-атаки, при которой
		// атакующий умышленно медленно отправляет HTTP-заголовки, удерживая соединения открытыми и истощая
		// пул доступных соединений на сервере. ReadHeaderTimeout принудительно закрывает соединение,
		// если клиент не успел отправить все заголовки за отведенное время.
	}

	// Запускаем сервер в отдельной горутине
	go func() {
		log.Printf("🚀 HTTP-сервер запущен на адресе %s\n", config.GetConfig().OrderHTTP.Address())
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
	for filepath.Base(dir) != "order" {
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", errors.New("order not found")
		}
		dir = parent
	}
	return dir, nil
}
