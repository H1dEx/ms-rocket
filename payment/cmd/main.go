package main

import (
	"context"
	"fmt"
	"net"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	paymentApi "github.com/H1dEx/ms-rocket/payment/internal/api/payment/v1"
	"github.com/H1dEx/ms-rocket/payment/internal/config"
	paymentService "github.com/H1dEx/ms-rocket/payment/internal/service/payment"
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

	defer func() {
		if cerr := lis.Close(); cerr != nil {
			logger.Error(appCtx, "failed to close listener", zap.Error(cerr))
		}
	}()

	s := grpc.NewServer()

	health.RegisterServer(s)

	service := paymentService.NewService()
	api := paymentApi.NewAPI(service)

	paymentV1.RegisterPaymentServiceServer(s, api)

	reflection.Register(s)

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
	s.GracefulStop()
	logger.Info(appCtx, "✅ Server stopped")
}
