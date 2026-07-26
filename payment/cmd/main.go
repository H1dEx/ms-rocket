package main

import (
	"context"
	"fmt"
	"net"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	paymentApi "github.com/H1dEx/ms-rocket/payment/internal/api/payment/v1"
	"github.com/H1dEx/ms-rocket/payment/internal/config"
	paymentService "github.com/H1dEx/ms-rocket/payment/internal/service/payment"
	"github.com/H1dEx/ms-rocket/platform/pkg/closer"
	"github.com/H1dEx/ms-rocket/platform/pkg/grpc/health"
	"github.com/H1dEx/ms-rocket/platform/pkg/logger"
	paymentV1 "github.com/H1dEx/ms-rocket/shared/pkg/proto/payment/v1"
)

const configPath = "./deploy/compose/payment/.env"

func main() {
	appCtx, appCancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer appCancel()

	err := logger.Init("info", true)
	if err != nil {
		fmt.Println("failed to init logger", err)
		return
	}

	closer.SetLogger(logger.Logger())
	defer gracefulShutdown()

	err = config.Load(configPath)
	if err != nil {
		logger.Error(appCtx, "failed to load config", zap.Error(err))
		return
	}
	lis, err := net.Listen("tcp", config.GetConfig().PaymentGRPC.Address())
	if err != nil {
		logger.Error(appCtx, "failed to listen", zap.Error(err))
		return
	}

	s := grpc.NewServer()

	health.RegisterServer(s)

	service := paymentService.NewService()
	api := paymentApi.NewAPI(service)

	paymentV1.RegisterPaymentServiceServer(s, api)

	reflection.Register(s)

	closer.AddNamed("gRPC server", func(ctx context.Context) error {
		return stopGRPCServer(ctx, s)
	})

	go func() {
		logger.Info(appCtx, fmt.Sprintf("🚀 gRPC server listening on %s", config.GetConfig().PaymentGRPC.Address()))
		err = s.Serve(lis)
		if err != nil {
			logger.Error(appCtx, "failed to serve", zap.Error(err))
			return
		}
	}()

	// Graceful shutdown
	<-appCtx.Done()
	logger.Info(appCtx, "🛑 Shutting down gRPC server...")
}

func stopGRPCServer(ctx context.Context, s *grpc.Server) error {
	done := make(chan struct{})
	go func() {
		s.GracefulStop()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		s.Stop()
		return ctx.Err()
	}
}

func gracefulShutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := closer.CloseAll(ctx); err != nil {
		logger.Error(ctx, "❌ Ошибка при завершении работы", zap.Error(err))
	}
}
