package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/H1dEx/ms-rocket/order/internal/config"
	"github.com/H1dEx/ms-rocket/platform/pkg/closer"
	"github.com/H1dEx/ms-rocket/platform/pkg/logger"
	orderV1 "github.com/H1dEx/ms-rocket/shared/pkg/openapi/order/v1"
)

const (
	readHeaderTimeout = 5 * time.Second
)

type App struct {
	diContainer *diContainer
	httpServer  *http.Server
}

func New(ctx context.Context) (*App, error) {
	a := &App{}
	if err := a.initDeps(ctx); err != nil {
		return nil, err
	}
	return a, nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initDI,
		a.InitLogger,
		a.initCloser,
		a.initHTTPServer,
	}

	for _, init := range inits {
		if err := init(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (a *App) Run(ctx context.Context) error {
	// Канал для ошибок от компонентов
	errCh := make(chan error, 2)

	// Контекст для остановки всех горутин
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		if err := a.runHTTPServer(ctx); err != nil {
			errCh <- errors.Errorf("HTTP server crashed: %v", err)
		}
	}()

	go func() {
		if err := a.runConsumer(ctx); err != nil {
			errCh <- errors.Errorf("Consumer crashed: %v", err)
		}
	}()

	select {
	case <-ctx.Done():
		logger.Info(ctx, "🛑 Consumer context cancelled, shutting down")
		return ctx.Err()
	case err := <-errCh:
		logger.Error(ctx, "💥 Consumer crashed", zap.Error(err))
		cancel()
		return err
	}
}

func (a *App) runConsumer(ctx context.Context) error {
	logger.Info(ctx, "Starting order consumer")

	err := a.diContainer.OrderConsumer(ctx).RunConsumer(ctx)
	if err != nil {
		logger.Error(ctx, "failed to run consumer", zap.Error(err))
		return err
	}
	return nil
}

func (a *App) initDI(_ context.Context) error {
	a.diContainer = NewDIContainer()
	return nil
}

func (a *App) initHTTPServer(ctx context.Context) error {
	orderServer, err := orderV1.NewServer(a.diContainer.OrderV1API(ctx))
	if err != nil {
		logger.Error(ctx, "failed to create server", zap.Error(err))
		return err
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))
	r.Mount("/", orderServer)

	server := &http.Server{
		Addr:              config.GetConfig().OrderHTTP.Address(),
		Handler:           r,
		ReadHeaderTimeout: readHeaderTimeout, // Защита от Slowloris атак - тип DDoS-атаки, при которой
		// атакующий умышленно медленно отправляет HTTP-заголовки, удерживая соединения открытыми и истощая
		// пул доступных соединений на сервере. ReadHeaderTimeout принудительно закрывает соединение,
		// если клиент не успел отправить все заголовки за отведенное время.
	}

	closer.AddNamed("HTTP server", func(ctx context.Context) error {
		return server.Shutdown(ctx)
	})
	a.httpServer = server
	return nil
}

func (a *App) runHTTPServer(ctx context.Context) error {
	logger.Info(ctx, fmt.Sprintf("Starting HTTP server listening on %s", config.GetConfig().OrderHTTP.Address()))

	err := a.httpServer.ListenAndServe()
	if err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			logger.Error(ctx, "failed to serve", zap.Error(err))
			return err
		}
	}
	return nil
}

func (a *App) InitLogger(_ context.Context) error {
	err := logger.Init(config.GetConfig().Logger.Level(), config.GetConfig().Logger.AsJSON())
	if err != nil {
		return err
	}
	return nil
}

func (a *App) initCloser(_ context.Context) error {
	closer.SetLogger(logger.Logger())
	return nil
}
