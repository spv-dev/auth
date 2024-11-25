package config

import (
	"errors"
	"net"
	"os"

	serviceerror "github.com/spv-dev/auth/internal/service_error"
)

const (
	httpHostEnvName = "HTTP_HOST"
	httpPortEnvName = "HTTP_PORT"
)

// HTTPConfig интерфейс для работы c HTTP сервером
type HTTPConfig interface {
	Address() string
}

type httpConfig struct {
	host string
	port string
}

// NewHTTPConfig получение конфигурации для подключения к HTTP серверу
func NewHTTPConfig() (HTTPConfig, error) {
	host := os.Getenv(httpHostEnvName)
	if len(host) == 0 {
		return nil, errors.New(serviceerror.HTTPHostNotFound)
	}

	port := os.Getenv(httpPortEnvName)
	if len(port) == 0 {
		return nil, errors.New(serviceerror.HTTPPortNotFound)
	}

	return &httpConfig{
		host: host,
		port: port,
	}, nil
}

func (cfg *httpConfig) Address() string {
	return net.JoinHostPort(cfg.host, cfg.port)
}
