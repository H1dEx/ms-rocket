package redis

import (
	"context"
	"time"

	redigo "github.com/gomodule/redigo/redis"
	"go.uber.org/zap"
)


type client struct {
	pool              *redigo.Pool
	logger            Logger
	connectionTimeout time.Duration
}

type Logger interface {
	Info(ctx context.Context, msg string, fields ...zap.Field)
	Error(ctx context.Context, msg string, fields ...zap.Field)
}

type redisFn func(ctx context.Context, conn redigo.Conn) error

func NewClient(pool *redigo.Pool, logger Logger, connectionTimeout time.Duration) *client {
	return &client{
		pool:              pool,
		logger:            logger,
		connectionTimeout: connectionTimeout,
	}
}

func (c *client) withConnection(ctx context.Context, fn redisFn) error {
	conn, err := c.getConntection(ctx)
	if err != nil {
		return err
	}
	defer func(){
		if err := conn.Close(); err != nil {
			c.logger.Error(ctx, "failed to close connection", zap.Error(err))
		}
	}()

	return fn(ctx, conn)
}

func (c *client) getConntection(ctx context.Context) (redigo.Conn, error) {
	ctx, cancel := context.WithTimeout(ctx, c.connectionTimeout)
	defer cancel()

	conn, err := c.pool.GetContext(ctx)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func (c *client) Set(ctx context.Context, key string, value any) error {
	return c.withConnection(ctx, func(ctx context.Context, conn redigo.Conn) error {
		_, err := conn.Do("SET", key, value)
		return err
	})
}

func (c *client) SetWithTTL(ctx context.Context, key string, value any, ttl time.Duration) error {
	return c.withConnection(ctx, func(ctx context.Context, conn redigo.Conn) error {
		_, err := conn.Do("SET", key, value, "EX", int(ttl.Seconds()))
		return err
	})
}

func (c *client) Get(ctx context.Context, key string) ([]byte, error) {
	var result []byte
	err := c.withConnection(ctx, func(ctx context.Context, conn redigo.Conn) error {
		val, err := redigo.Bytes(conn.Do("GET", key))
		if err != nil {
			return err
		}
		result = val
		return nil
	})
	return result, err
}

func (c *client) HashSet(ctx context.Context, key string, values any) error {
	return c.withConnection(ctx, func(ctx context.Context, conn redigo.Conn) error {
		_, err := conn.Do("HSET", redigo.Args{key}.AddFlat(values)...)
		return err
	})
}

func (c *client) HGetAll(ctx context.Context, key string) ([]any, error) {
	var result []any
	err := c.withConnection(ctx, func(ctx context.Context, conn redigo.Conn) error {
		val, err :=redigo.Values(conn.Do("HGETALL", key))
		if err != nil {
			return err
		}
		result = val
		return nil
	})
	return result, err
}

func (c *client) Del (ctx context.Context, key string) error {
	return c.withConnection(ctx, func(ctx context.Context, conn redigo.Conn) error {
		_, err := conn.Do("DEL", key)
		return err
	})
}

func (c *client) Exists(ctx context.Context, key string) (bool, error) {
	var result bool
	err := c.withConnection(ctx, func(ctx context.Context, conn redigo.Conn)error{
		val, err := redigo.Bool(conn.Do("EXISTS", key))
		if err != nil {
			return err
		}
		result = val
		return nil
	})
	return result, err
}

func (c *client) Ping(ctx context.Context) error {
	return c.withConnection(ctx, func(ctx context.Context, conn redigo.Conn) error {
		_, err := conn.Do("PING")
		return err
	})
}

func (c client) Expire(ctx context.Context, key string, ttl time.Duration) error {
	return c.withConnection(ctx, func(ctx context.Context, conn redigo.Conn) error {
		_, err := conn.Do("EXPIRE", key, int(ttl.Seconds()))
		return err
	})
}