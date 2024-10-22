package config

import (
	"errors"
	"os"

	"github.com/spv-dev/auth/internal/config"
)

const (
	dsnEnvName = "PG_DSN"
)

var _ config.PGConfig = (*pgConfig)(nil)

type pgConfig struct {
	dsn string
}

// NewPGConfig получает конфигурацию для Postgres
func NewPGConfig() (*pgConfig, error) {
	dsn := os.Getenv(dsnEnvName)
	if len(dsn) == 0 {
		return nil, errors.New("pg dsn not found")
	}

	return &pgConfig{
		dsn: dsn,
	}, nil
}

// DSN возвращает dsn конфигурации
func (cfg *pgConfig) DSN() string {
	return cfg.dsn
}
