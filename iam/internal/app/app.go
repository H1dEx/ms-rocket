package app

import (
	"context"
	"fmt"
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/H1dEx/ms-rocket/iam/internal/config"
	"github.com/H1dEx/ms-rocket/platform/pkg/closer"
	"github.com/H1dEx/ms-rocket/platform/pkg/grpc/health"
	"github.com/H1dEx/ms-rocket/platform/pkg/logger"
	authV1 "github.com/H1dEx/ms-rocket/shared/pkg/proto/auth/v1"
	userV1 "github.com/H1dEx/ms-rocket/shared/pkg/proto/user/v1"
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

func (a *App) initLogger(ctx context.Context) error {
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
	listener, err := net.Listen("tcp", config.GetConfig().IAMGRPC.Address())
	if err != nil {
		return err
	}
	a.listener = listener
	return nil
}

func (a *App) initGRPCServer(ctx context.Context) error {
	userAPI := a.diContainer.UserV1API(ctx)
	authAPI := a.diContainer.AuthV1API(ctx)
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
		}
	})
	reflection.Register(a.grpcServer)
	health.RegisterServer(a.grpcServer)
	userV1.RegisterUserServiceServer(a.grpcServer, userAPI)
	authV1.RegisterAuthServiceServer(a.grpcServer, authAPI)
	return nil
}

func (a *App) runGRPCServer(ctx context.Context) error {
	logger.Info(ctx, fmt.Sprintf("Starting gRPC server listening on %s", config.GetConfig().IAMGRPC.Address()))
	err := a.grpcServer.Serve(a.listener)
	if err != nil {
		logger.Error(ctx, "failed to serve", zap.Error(err))
		return err
	}
	return nil
}
