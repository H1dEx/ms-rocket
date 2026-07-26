package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderApi "github.com/H1dEx/ms-rocket/order/internal/api/order/v1"
	inventoryCli "github.com/H1dEx/ms-rocket/order/internal/client/grpc/inventory/v1"
	paymentCli "github.com/H1dEx/ms-rocket/order/internal/client/grpc/payment/v1"
	"github.com/H1dEx/ms-rocket/order/internal/config"
	"github.com/H1dEx/ms-rocket/order/internal/migrator"
	orderRepo "github.com/H1dEx/ms-rocket/order/internal/repository/order"
	orderService "github.com/H1dEx/ms-rocket/order/internal/service/order"
	"github.com/H1dEx/ms-rocket/platform/pkg/logger"
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

	orderDir, err := findOrderDir()
	if err != nil {
		logger.Error(appCtx, "failed to find order directory", zap.Error(err))
		return
	}

	ctx := context.Background()
	conn, err := pgxpool.New(ctx, config.GetConfig().Postgres.URI())
	if err != nil {
		logger.Error(appCtx, "failed to connect to database", zap.Error(err))
		return
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err = conn.Ping(ctx)
	if err != nil {
		logger.Error(appCtx, "failed to ping database", zap.Error(err))
		return
	}
	sqlDB := stdlib.OpenDB(*conn.Config().ConnConfig)
	defer func() {
		if cerr := sqlDB.Close(); cerr != nil {
			logger.Error(appCtx, "failed to close database", zap.Error(cerr))
		}
	}()
	migrationDir := filepath.Join(orderDir, config.GetConfig().Postgres.MigrationDir())
	migrator := migrator.NewMigrator(sqlDB, migrationDir)
	err = migrator.Up()
	if err != nil {
		logger.Error(appCtx, "failed to migrate database", zap.Error(err))
		return
	}

	paymentConn, err := grpc.NewClient(
		config.GetConfig().PaymentGRPC.Address(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		logger.Error(appCtx, "failed to connect", zap.Error(err))
		return
	}
	defer func() {
		if cerr := paymentConn.Close(); cerr != nil {
			logger.Error(appCtx, "failed to close connect", zap.Error(cerr))
		}
	}()

	paymentClient := paymentV1.NewPaymentServiceClient(paymentConn)

	inventoryConn, err := grpc.NewClient(
		config.GetConfig().InventoryGRPC.Address(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		logger.Error(appCtx, "failed to connect", zap.Error(err))
		return
	}
	defer func() {
		if cerr := inventoryConn.Close(); cerr != nil {
			logger.Error(appCtx, "failed to close connect", zap.Error(cerr))
		}
	}()

	inventoryClient := inventoryV1.NewInventoryServiceClient(inventoryConn)
	repo := orderRepo.NewOrderRepository(conn)
	service := orderService.NewOrderService(repo, inventoryCli.NewInventoryClient(inventoryClient), paymentCli.NewPaymentClient(paymentClient))
	api := orderApi.NewOrderAPI(service)

	orderServer, err := orderV1.NewServer(api)
	if err != nil {
		logger.Error(appCtx, "failed to create server", zap.Error(err))
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
		logger.Info(appCtx, fmt.Sprintf("🚀 HTTP-сервер запущен на адресе %s\n", config.GetConfig().OrderHTTP.Address()))
		err = server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error(appCtx, "failed to start server", zap.Error(err))
		}
	}()

	// Graceful shutdown
	<-appCtx.Done()

	logger.Info(appCtx, "🛑 Shutting down HTTP server...")

	// Создаем контекст с таймаутом для остановки сервера
	ctx, cancel = context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		logger.Error(appCtx, "failed to stop server", zap.Error(err))
	}

	logger.Info(appCtx, "✅ Server stopped")
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
