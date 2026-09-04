package app

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"runtime"

	redigo "github.com/gomodule/redigo/redis"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"

	authApi "github.com/H1dEx/ms-rocket/iam/internal/api/auth/v1"
	userApi "github.com/H1dEx/ms-rocket/iam/internal/api/user/v1"
	"github.com/H1dEx/ms-rocket/iam/internal/config"
	"github.com/H1dEx/ms-rocket/iam/internal/repository"
	sessionRepo "github.com/H1dEx/ms-rocket/iam/internal/repository/session"
	userRepo "github.com/H1dEx/ms-rocket/iam/internal/repository/user"
	"github.com/H1dEx/ms-rocket/iam/internal/service"
	authService "github.com/H1dEx/ms-rocket/iam/internal/service/auth"
	userService "github.com/H1dEx/ms-rocket/iam/internal/service/user"
	"github.com/H1dEx/ms-rocket/platform/pkg/cache"
	"github.com/H1dEx/ms-rocket/platform/pkg/cache/redis"
	"github.com/H1dEx/ms-rocket/platform/pkg/closer"
	"github.com/H1dEx/ms-rocket/platform/pkg/logger"
	"github.com/H1dEx/ms-rocket/platform/pkg/migrator"
	authV1 "github.com/H1dEx/ms-rocket/shared/pkg/proto/auth/v1"
	userV1 "github.com/H1dEx/ms-rocket/shared/pkg/proto/user/v1"
)

type diContainer struct {
	userV1API   userV1.UserServiceServer
	userService service.UserService
	userRepo    repository.UserRepository

	authV1API   authV1.AuthServiceServer
	authService service.AuthService
	sessionRepo repository.SessionRepository

	redisPool   *redigo.Pool
	redisClient cache.RedisClient

	postgresConn *pgxpool.Pool
}

func NewDIContainer() *diContainer {
	return &diContainer{}
}

func (c *diContainer) UserV1API(ctx context.Context) userV1.UserServiceServer {
	if c.userV1API == nil {
		c.userV1API = userApi.NewUserAPI(c.UserService(ctx))
	}
	return c.userV1API
}

func (c *diContainer) UserService(ctx context.Context) service.UserService {
	if c.userService == nil {
		c.userService = userService.NewService(c.UserRepo(ctx))
	}
	return c.userService
}

func (c *diContainer) UserRepo(ctx context.Context) repository.UserRepository {
	if c.userRepo == nil {
		c.userRepo = userRepo.NewRepository(c.PostgresConn(ctx))
	}
	return c.userRepo
}

func (c *diContainer) AuthV1API(ctx context.Context) authV1.AuthServiceServer {
	if c.authV1API == nil {
		c.authV1API = authApi.NewAuthApi(c.AuthService(ctx))
	}
	return c.authV1API
}

func (c *diContainer) AuthService(ctx context.Context) service.AuthService {
	if c.authService == nil {
		c.authService = authService.NewService(c.SessionRepo(ctx), c.UserService(ctx), config.GetConfig().Session.TTL())
	}
	return c.authService
}

func (c *diContainer) SessionRepo(ctx context.Context) repository.SessionRepository {
	if c.sessionRepo == nil {
		c.sessionRepo = sessionRepo.NewRepository(c.RedisClient(ctx))
	}

	return c.sessionRepo
}

func (c *diContainer) RedisPool(_ context.Context) *redigo.Pool {
	if c.redisPool == nil {
		c.redisPool = &redigo.Pool{
			MaxIdle:     config.GetConfig().Redis.MaxIdle(),
			IdleTimeout: config.GetConfig().Redis.IdleTimeout(),
			DialContext: func(ctx context.Context) (redigo.Conn, error) {
				return redigo.DialContext(ctx, "tcp", config.GetConfig().Redis.Address())
			},
		}
		closer.AddNamed("Redis pool", func(_ context.Context) error {
			return c.redisPool.Close()
		})
	}
	return c.redisPool
}

func (c *diContainer) RedisClient(ctx context.Context) cache.RedisClient {
	if c.redisClient == nil {
		c.redisClient = redis.NewClient(c.RedisPool(ctx), logger.Logger(), config.GetConfig().Redis.ConnectionTimeout())
		if err := c.redisClient.Ping(ctx); err != nil {
			panic(fmt.Errorf("failed to ping Redis: %w", err))
		}
	}
	return c.redisClient
}

func findIAMDir() (string, error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", errors.New("runtime.Caller failed")
	}
	dir := filepath.Dir(file)
	for filepath.Base(dir) != "iam" {
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", errors.New("iam not found")
		}
		dir = parent
	}
	return dir, nil
}

func (c *diContainer) PostgresConn(ctx context.Context) *pgxpool.Pool {
	if c.postgresConn == nil {
		conn, err := pgxpool.New(ctx, config.GetConfig().Postgres.URI())
		if err != nil {
			panic(fmt.Errorf("failed to connect to database: %s", err.Error()))
		}
		closer.AddNamed("PostgreSQL connection", func(_ context.Context) error {
			conn.Close()
			return nil
		})
		err = conn.Ping(ctx)
		if err != nil {
			panic(fmt.Errorf("failed to ping database: %s", err.Error()))
		}
		sqlDB := stdlib.OpenDB(*conn.Config().ConnConfig)
		closer.AddNamed("PostgreSQL converted connection", func(ctx context.Context) error {
			if cerr := sqlDB.Close(); cerr != nil {
				logger.Error(ctx, "failed to close database", zap.Error(cerr))
			}
			return nil
		})

		iamDir, err := findIAMDir()
		if err != nil {
			panic(fmt.Errorf("failed to find iam directory: %s", err.Error()))
		}
		migrationDir := filepath.Join(iamDir, "..", config.GetConfig().Postgres.MigrationDir())
		migrator := migrator.NewMigrator(sqlDB, migrationDir)
		err = migrator.Up()
		if err != nil {
			panic(fmt.Errorf("failed to migrate database: %s", err.Error()))
		}
		c.postgresConn = conn
	}
	return c.postgresConn
}
