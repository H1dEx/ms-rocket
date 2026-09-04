package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/H1dEx/ms-rocket/iam/internal/app"
	"github.com/H1dEx/ms-rocket/iam/internal/config"
	"github.com/H1dEx/ms-rocket/platform/pkg/closer"
	"github.com/H1dEx/ms-rocket/platform/pkg/logger"
)

const configPath = "./deploy/compose/iam/.env"

func main() {
	err := config.Load(configPath)
	if err != nil {
		panic(fmt.Errorf("failed to load config: %w", err))
	}
	appCtx, appCancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer appCancel()
	defer gracefulShutdown()

	app, err := app.New(appCtx)
	if err != nil {
		logger.Error(appCtx, "failed to create app", zap.Error(err))
		return
	}

	go func() {
		if err := app.Run(appCtx); err != nil {
			logger.Error(appCtx, "failed to run app", zap.Error(err))
			appCancel()
			return
		}
	}()
	<-appCtx.Done()
}

func gracefulShutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := closer.CloseAll(ctx); err != nil {
		logger.Error(ctx, "❌ Ошибка при завершении работы", zap.Error(err))
	}
}
