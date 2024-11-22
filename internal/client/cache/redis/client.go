package redis

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/spv-dev/auth/internal/client/cache"
	"github.com/spv-dev/auth/internal/config"
)

var _ cache.RedisClient = (*client)(nil)

type handler func(ctx context.Context, conn redis.Conn) error

type client struct {
	pool   *redis.Pool
	config config.RedisConfig
}

// NewClient создание нового клиента для Redis
func NewClient(pool *redis.Pool, config config.RedisConfig) *client {
	return &client{
		pool:   pool,
		config: config,
	}
}

// HashSet обёртка для выполнения HSET в Redis. Сохраняет объект
func (c *client) HashSet(ctx context.Context, key string, values interface{}) error {
	err := c.execute(ctx, func(_ context.Context, conn redis.Conn) error {
		_, err := conn.Do("HSET", redis.Args{key}.AddFlat(values)...)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

// Set обёртка для выполнения SET в Redis. Сохраняет одно значение
func (c *client) Set(ctx context.Context, key string, value interface{}) error {
	err := c.execute(ctx, func(_ context.Context, conn redis.Conn) error {
		_, err := conn.Do("SET", redis.Args{key}.Add(value)...)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

// HGetAll обёртка для выполнения SET в Redis. Получает объект
func (c *client) HGetAll(ctx context.Context, key string) ([]interface{}, error) {
	var values []interface{}
	err := c.execute(ctx, func(_ context.Context, conn redis.Conn) error {
		var errEx error
		values, errEx = redis.Values(conn.Do("HGETALL", key))
		if errEx != nil {
			return errEx
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return values, nil
}

// Get обёртка для выполнения SET в Redis. Получает одно значение
func (c *client) Get(ctx context.Context, key string) (interface{}, error) {
	var value interface{}
	err := c.execute(ctx, func(_ context.Context, conn redis.Conn) error {
		var errEx error
		value, errEx = conn.Do("GET", key)
		if errEx != nil {
			return errEx
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return value, nil
}

// Expire обёртка для выполнения EXPIRE в Redis. Установка времени жизни
func (c *client) Expire(ctx context.Context, key string, expiration time.Duration) error {
	err := c.execute(ctx, func(_ context.Context, conn redis.Conn) error {
		_, err := conn.Do("EXPIRE", key, int(expiration.Seconds()))
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

// Ping обёртка для выполнения PING в Redis. Проверка связи с хранилищем
func (c *client) Ping(ctx context.Context) error {
	err := c.execute(ctx, func(_ context.Context, conn redis.Conn) error {
		_, err := conn.Do("PING")
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *client) execute(ctx context.Context, handler handler) error {
	conn, err := c.getConnect(ctx)
	if err != nil {
		return err
	}

	defer func() {
		err = conn.Close()
		if err != nil {
			log.Printf("failed to close redis connection: %v", err)
		}
	}()

	err = handler(ctx, conn)
	if err != nil {
		return err
	}

	return nil
}

func (c *client) getConnect(ctx context.Context) (redis.Conn, error) {
	getConnectTimoutCtx, cancel := context.WithTimeout(ctx, c.config.ConnectionTimeout())
	defer cancel()

	conn, err := c.pool.GetContext(getConnectTimoutCtx)
	if err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("failed to get redis connection: %v", err)
	}

	return conn, nil
}
