package config

import (
	"github.com/joho/godotenv"

	"github.com/H1dEx/ms-rocket/iam/internal/config/env"
)

var appConfig *config

type config struct {
	Logger   LoggerConfig
	IAMGRPC  IAMGRPCConfig
	Postgres PostgresConfig
	Redis    RedisConfig
	Session  SessionConfig
}

func Load(path ...string) error {
	err := godotenv.Load(path...)
	if err != nil {
		return err
	}
	loggerConfig, err := env.NewLoggerConfig()
	if err != nil {
		return err
	}
	iamGRPCConfig, err := env.NewIAMGRPCConfig()
	if err != nil {
		return err
	}
	postgresConfig, err := env.NewPostgresConfig()
	if err != nil {
		return err
	}
	redisConfig, err := env.NewRedisConfig()
	if err != nil {
		return err
	}
	sessionConfig, err := env.NewSessionConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger:   loggerConfig,
		IAMGRPC:  iamGRPCConfig,
		Postgres: postgresConfig,
		Redis:    redisConfig,
		Session:  sessionConfig,
	}
	return nil
}

func GetConfig() *config {
	return appConfig
}
