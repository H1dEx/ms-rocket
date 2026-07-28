package app

import (
	"context"
	"fmt"
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/H1dEx/ms-rocket/inventory/internal/config"
	"github.com/H1dEx/ms-rocket/platform/pkg/closer"
	"github.com/H1dEx/ms-rocket/platform/pkg/grpc/health"
	"github.com/H1dEx/ms-rocket/platform/pkg/logger"
	inventoryV1 "github.com/H1dEx/ms-rocket/shared/pkg/proto/inventory/v1"
)

type App struct {
	diContainer *diContainer
	grpcServer  *grpc.Server
	listener    net.Listener
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
		a.initLogger,
		a.initCloser,
		a.initListener,
		a.initGRPCServer,
	}

	for _, init := range inits {
		if err := init(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (a *App) Run(ctx context.Context) error {
	return a.runGRPCServer(ctx)
}

func (a *App) initDI(_ context.Context) error {
	a.diContainer = NewDIContainer()
	return nil
}

func (a *App) initLogger(_ context.Context) error {
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

func (a *App) initListener(_ context.Context) error {
	lis, err := net.Listen("tcp", config.GetConfig().InventoryGRPC.Address())
	if err != nil {
		return err
	}

	a.listener = lis
	return nil
}

func (a *App) initGRPCServer(ctx context.Context) error {
	api := a.diContainer.InventoryV1API(ctx)
	a.grpcServer = grpc.NewServer()
	closer.AddNamed("gRPC server", func(ctx context.Context) error {
		done := make(chan struct{})
		go func() {
			a.grpcServer.GracefulStop()
			close(done)
		}()
		select {
		case <-done:
			return nil
		case <-ctx.Done():
			a.grpcServer.Stop()
			return ctx.Err()
	}})

	reflection.Register(a.grpcServer)

	health.RegisterServer(a.grpcServer)

	inventoryV1.RegisterInventoryServiceServer(a.grpcServer, api)

	return nil
}

func (a *App) runGRPCServer(ctx context.Context) error {
	logger.Info(ctx, fmt.Sprintf("Starting gRPC server listening on %s", config.GetConfig().InventoryGRPC.Address()))

	err := a.grpcServer.Serve(a.listener)
	if err != nil {
		logger.Error(ctx, "failed to serve", zap.Error(err))
		return err
	}
	return nil
}
