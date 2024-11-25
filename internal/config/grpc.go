package config

import (
	"errors"
	"net"
	"os"

	serviceerror "github.com/spv-dev/auth/internal/service_error"
)

const (
	grpcHostEnvName = "GRPC_HOST"
	grpcPortEnvName = "GRPC_PORT"
)

// GRPCConfig интерфейс для конфигурации gRPC
type GRPCConfig interface {
	Address() string
}

type grpcConfig struct {
	host string
	port string
}

// NewGRPCConfig получает конфигурацию gRPC
func NewGRPCConfig() (*grpcConfig, error) {
	host := os.Getenv(grpcHostEnvName)
	if len(host) == 0 {
		return nil, errors.New(serviceerror.GRPCHostNotFound)
	}

	port := os.Getenv(grpcPortEnvName)
	if len(port) == 0 {
		return nil, errors.New(serviceerror.GRPCPortNotFound)
	}

	return &grpcConfig{
		host: host,
		port: port,
	}, nil
}

// Address возвращает адрес сервиса
func (cfg *grpcConfig) Address() string {
	return net.JoinHostPort(cfg.host, cfg.port)
}
