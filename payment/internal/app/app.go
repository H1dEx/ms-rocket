package app

import (
	"context"
	"fmt"
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/H1dEx/ms-rocket/payment/internal/config"
	"github.com/H1dEx/ms-rocket/platform/pkg/closer"
	"github.com/H1dEx/ms-rocket/platform/pkg/grpc/health"
	"github.com/H1dEx/ms-rocket/platform/pkg/logger"
	paymentV1 "github.com/H1dEx/ms-rocket/shared/pkg/proto/payment/v1"
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
		a.InitLogger,
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

func (a *App) initListener(_ context.Context) error {
	lis, err := net.Listen("tcp", config.GetConfig().PaymentGRPC.Address())
	if err != nil {
		return err
	}

	a.listener = lis
	return nil
}

func (a *App) initGRPCServer(_ context.Context) error {
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

	paymentV1.RegisterPaymentServiceServer(a.grpcServer, a.diContainer.PaymentV1API())

	return nil
}

func (a *App) runGRPCServer(ctx context.Context) error {
	logger.Info(ctx, fmt.Sprintf("Starting gRPC server listening on %s", config.GetConfig().PaymentGRPC.Address()))

	err := a.grpcServer.Serve(a.listener)
	if err != nil {
		logger.Error(ctx, "failed to serve", zap.Error(err))
		return err
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
