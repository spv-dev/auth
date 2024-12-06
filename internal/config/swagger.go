package config

import (
	"errors"
	"net"
	"os"

	serviceerror "github.com/spv-dev/auth/internal/service_error"
)

const (
	swaggerHostEnvName = "SWAGGER_HOST"
	swaggerPortEnvName = "SWAGGER_PORT"
)

// SwaggerConfig интерфейс для получения информации о Swagger
type SwaggerConfig interface {
	Address() string
}

type swaggerConfig struct {
	host string
	port string
}

// NewSwaggerConfig получение конфигурации Swagger
func NewSwaggerConfig() (SwaggerConfig, error) {
	host := os.Getenv(swaggerHostEnvName)
	if len(host) == 0 {
		return nil, errors.New(serviceerror.SwaggerHostNotFound)
	}

	port := os.Getenv(swaggerPortEnvName)
	if len(port) == 0 {
		return nil, errors.New(serviceerror.SwaggerPortNotFound)
	}

	return &swaggerConfig{
		host: host,
		port: port,
	}, nil
}

func (cfg *swaggerConfig) Address() string {
	return net.JoinHostPort(cfg.host, cfg.port)
}
