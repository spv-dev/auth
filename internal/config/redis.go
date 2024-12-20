package config

import (
	"errors"
	"net"
	"os"
	"strconv"
	"time"

	serviceerror "github.com/spv-dev/auth/internal/service_error"
)

const (
	redisHostEnvName              = "REDIS_HOST"
	redisPortEnvName              = "REDIS_PORT"
	redisConnectionTimeoutEnvName = "REDIS_CONNECTION_TIMEOUT_SEC"
	redisMaxIdleEnvName           = "REDIS_MAX_IDLE"
	redisIdleTimeoutEnvName       = "REDIS_IDLE_TIMEOUT_SEC"
)

// RedisConfig интерфейс конфигурации Redis
type RedisConfig interface {
	Address() string
	ConnectionTimeout() time.Duration
	MaxIdle() int
	IdleTimeout() time.Duration
}

type redisConfig struct {
	host string
	port string

	connectionTimeout time.Duration

	maxIdle     int
	idleTimeout time.Duration
}

// NewRedisConfig получение конфигурации для подсоединение к Redis
func NewRedisConfig() (*redisConfig, error) {
	host := os.Getenv(redisHostEnvName)
	if len(host) == 0 {
		return nil, errors.New(serviceerror.RedisHostNotFound)
	}

	port := os.Getenv(redisPortEnvName)
	if len(port) == 0 {
		return nil, errors.New(serviceerror.RedisPortNotFound)
	}

	connectionTimeoutStr := os.Getenv(redisConnectionTimeoutEnvName)
	if len(connectionTimeoutStr) == 0 {
		return nil, errors.New(serviceerror.RedisConnectionTimeoutNotFound)
	}

	connectionTimeout, err := strconv.ParseInt(connectionTimeoutStr, 10, 64)
	if err != nil {
		return nil, errors.New(serviceerror.FailedToParseConnectionTimeout)
	}

	maxIdleStr := os.Getenv(redisMaxIdleEnvName)
	if len(connectionTimeoutStr) == 0 {
		return nil, errors.New(serviceerror.RedisMaxIdleNotFound)
	}

	maxIdle, err := strconv.Atoi(maxIdleStr)
	if err != nil {
		return nil, errors.New(serviceerror.FailedToParseMaxIdle)
	}

	idleTimeoutStr := os.Getenv(redisIdleTimeoutEnvName)
	if len(idleTimeoutStr) == 0 {
		return nil, errors.New(serviceerror.RedisIdleTimeoutNotFound)
	}

	idleTimeout, err := strconv.ParseInt(idleTimeoutStr, 10, 64)
	if err != nil {
		return nil, errors.New(serviceerror.FailedToParseIdleTimeout)
	}

	return &redisConfig{
		host:              host,
		port:              port,
		connectionTimeout: time.Duration(connectionTimeout) * time.Second,
		maxIdle:           maxIdle,
		idleTimeout:       time.Duration(idleTimeout) * time.Second,
	}, nil
}

// Address получение адреса соединения с Redis
func (cfg *redisConfig) Address() string {
	return net.JoinHostPort(cfg.host, cfg.port)
}

// ConnectionTimeout получение значения времени соединения с Redis
func (cfg *redisConfig) ConnectionTimeout() time.Duration {
	return cfg.connectionTimeout
}

// MaxIdle получение значения максимального значения соединений с Redis
func (cfg *redisConfig) MaxIdle() int {
	return cfg.maxIdle
}

// IdleTimeout получение значения времени жизни Idle Redis
func (cfg *redisConfig) IdleTimeout() time.Duration {
	return cfg.idleTimeout
}
