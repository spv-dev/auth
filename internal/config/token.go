package config

import (
	"errors"
	"os"
	"strconv"
	"time"
)

const (
	refreshTokenSecretKeyEnvName  = "REFRESH_TOKEN_SECRET"
	refreshTokenExpirationEnvName = "REFRESH_TOKEN_EXPIRATION"
	accessTokenSecretKeyEnvName   = "ACCESS_TOKEN_SECRET"
	accessTokenExpirationEnvName  = "ACCESS_TOKEN_EXPIRATION"
)

type TokenConfig interface {
	GetRefreshSecret() string
	GetRefreshExpiration() time.Duration
	GetAccessSecret() string
	GetAccessExpiration() time.Duration
}

type tokenConfig struct {
	refreshTokenSecretKey  string
	refreshTokenExpiration time.Duration
	accessTokenSecretKey   string
	accessTokenExpiration  time.Duration
}

func NewTokenConfig() (*tokenConfig, error) {
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

	return &tokenConfig{
		refreshTokenSecretKey:  refreshTokenSecretKey,
		refreshTokenExpiration: time.Duration(refreshTokenExpiration * int64(time.Minute)),
		accessTokenSecretKey:   accessTokenSecretKey,
		accessTokenExpiration:  time.Duration(accessTokenExpiration * int64(time.Minute)),
	}, nil
}

func (cfg *tokenConfig) GetRefreshSecret() string {
	return cfg.refreshTokenSecretKey
}

func (cfg *tokenConfig) GetRefreshExpiration() time.Duration {
	return cfg.refreshTokenExpiration
}

func (cfg *tokenConfig) GetAccessSecret() string {
	return cfg.accessTokenSecretKey
}

func (cfg *tokenConfig) GetAccessExpiration() time.Duration {
	return cfg.accessTokenExpiration
}
