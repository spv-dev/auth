package config

import (
	"errors"
	"net"
	"os"

	"github.com/spv-dev/auth/internal/config"
)

const (
	grpcHostEnvName = "GRPC_HOST"
	grpcPortEnvName = "GRPC_PORT"
)

var _ config.GRPCConfig = (*grpcConfig)(nil)

type grpcConfig struct {
	host string
	port string
}

// NewGRPCConfig получает конфигурацию gRPC
func NewGRPCConfig() (*grpcConfig, error) {
	host := os.Getenv(grpcHostEnvName)
	if len(host) == 0 {
		return nil, errors.New("gRPC host not found")
	}

	port := os.Getenv(grpcPortEnvName)
	if len(port) == 0 {
		return nil, errors.New("gRPC port not found")
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
