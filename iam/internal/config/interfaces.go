package config

import "time"

type LoggerConfig interface {
	Level() string
	AsJSON() bool
}

type IAMGRPCConfig interface {
	Address() string
}

type PostgresConfig interface {
	URI() string
	MigrationDir() string
}

type RedisConfig interface {
	Address() string
	ConnectionTimeout() time.Duration
	MaxIdle() int
	IdleTimeout() time.Duration
}

type SessionConfig interface {
	TTL() time.Duration
}
