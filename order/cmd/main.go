package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/H1dEx/ms-rocket/order/internal/app"
	"github.com/H1dEx/ms-rocket/order/internal/config"
	"github.com/H1dEx/ms-rocket/platform/pkg/closer"
	"github.com/H1dEx/ms-rocket/platform/pkg/logger"
)

const (
	shutdownTimeout = 10 * time.Second
	configPath      = "./deploy/compose/order/.env"
)

func main() {
	err := config.Load(configPath)
	if err != nil {
		panic(fmt.Errorf("failed to load config: %w", err))
	}

	appCtx, appCancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer appCancel()
	defer gracefulShutdown()

	a, err := app.New(appCtx)
	if err != nil {
		logger.Error(appCtx, "failed to create app", zap.Error(err))
		return
	}
	go func() {
		if err := a.Run(appCtx); err != nil {
			logger.Error(appCtx, "failed to run app", zap.Error(err))
			appCancel()
			return
		}
	}()
	<-appCtx.Done()
}

func gracefulShutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := closer.CloseAll(ctx); err != nil {
		logger.Error(ctx, "❌ Ошибка при завершении работы", zap.Error(err))
	}
}
