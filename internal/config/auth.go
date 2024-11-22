package config

import (
	"errors"
	"net"
	"os"
	"strconv"
	"time"
)

const (
	authHostEnvName               = "AUTH_HOST"
	authPortEnvName               = "AUTH_PORT"
	refreshTokenSecretKeyEnvName  = "REFRESH_TOKEN_SECRET"     // #nosec G101
	refreshTokenExpirationEnvName = "REFRESH_TOKEN_EXPIRATION" // #nosec G101
	accessTokenSecretKeyEnvName   = "ACCESS_TOKEN_SECRET"      // #nosec G101
	accessTokenExpirationEnvName  = "ACCESS_TOKEN_EXPIRATION"  // #nosec G101
)

// AuthConfig интерфейс для конфигурации сервиса авторизации
type AuthConfig interface {
	Address() string
	GetRefreshSecret() string
	GetRefreshExpiration() time.Duration
	GetAccessSecret() string
	GetAccessExpiration() time.Duration
}

type authConfig struct {
	host                   string
	port                   string
	refreshTokenSecretKey  string
	refreshTokenExpiration time.Duration
	accessTokenSecretKey   string
	accessTokenExpiration  time.Duration
}

// NewAuthConfig получение новой конфигурации для сервиса авторизации
func NewAuthConfig() (*authConfig, error) {
	host := os.Getenv(authHostEnvName)
	if len(host) == 0 {
		return nil, errors.New("auth host not found")
	}

	port := os.Getenv(authPortEnvName)
	if len(port) == 0 {
		return nil, errors.New("auth port not found")
	}

	refreshTokenSecretKey := os.Getenv(refreshTokenSecretKeyEnvName)
	if len(refreshTokenSecretKey) == 0 {
		return nil, errors.New("refresh token secret not found")
	}

	refreshTokenExpirationStr := os.Getenv(refreshTokenExpirationEnvName)
	if len(refreshTokenExpirationStr) == 0 {
		return nil, errors.New("refresh token expiration not found")
	}

	refreshTokenExpiration, err := strconv.ParseInt(refreshTokenExpirationStr, 10, 64)
	if err != nil {
		return nil, errors.New("failed to parse refresh token expiration")
	}

	accessTokenSecretKey := os.Getenv(accessTokenSecretKeyEnvName)
	if len(accessTokenSecretKey) == 0 {
		return nil, errors.New("access token secret not found")
	}

	accessTokenExpirationStr := os.Getenv(accessTokenExpirationEnvName)
	if len(accessTokenExpirationStr) == 0 {
		return nil, errors.New("access token expiration not found")
	}

	accessTokenExpiration, err := strconv.ParseInt(accessTokenExpirationStr, 10, 64)
	if err != nil {
		return nil, errors.New("failed to parse connection timeout")
	}

	return &authConfig{
		host:                   host,
		port:                   port,
		refreshTokenSecretKey:  refreshTokenSecretKey,
		refreshTokenExpiration: time.Duration(refreshTokenExpiration * int64(time.Minute)),
		accessTokenSecretKey:   accessTokenSecretKey,
		accessTokenExpiration:  time.Duration(accessTokenExpiration * int64(time.Minute)),
	}, nil
}

func (cfg *authConfig) GetRefreshSecret() string {
	return cfg.refreshTokenSecretKey
}

func (cfg *authConfig) GetRefreshExpiration() time.Duration {
	return cfg.refreshTokenExpiration
}

func (cfg *authConfig) GetAccessSecret() string {
	return cfg.accessTokenSecretKey
}

func (cfg *authConfig) GetAccessExpiration() time.Duration {
	return cfg.accessTokenExpiration
}

// Address возвращает адрес сервиса
func (cfg *authConfig) Address() string {
	return net.JoinHostPort(cfg.host, cfg.port)
}
