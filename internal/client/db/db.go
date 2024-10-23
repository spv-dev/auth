package db

import (
	"context"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

// Client интерфейс для работы с БД через обёртки
type Client interface {
	DB() DB
	Close() error
}

// Query структура запроса
// Name - название запроса (чтобы можно было понятнее описать место вызова)
// QueryRaw - сам запрос
type Query struct {
	Name     string
	QueryRaw string
}

// SQLExecer интерфейс объединяет обёртки для работы с БД
type SQLExecer interface {
	NamedExecer
	QueryExecer
}

// NamedExecer интерфейс-обёртка для выполнения запросов в БД и возвращения значений
type NamedExecer interface {
	ScanOneContext(ctx context.Context, dest interface{}, q Query, args ...interface{}) error
	ScanAllContext(ctx context.Context, dest interface{}, q Query, args ...interface{}) error
}

// QueryExecer интерфейс-обёртка для выполнения запросов в БД
type QueryExecer interface {
	ExecContext(ctx context.Context, q Query, args ...interface{}) (pgconn.CommandTag, error)
	QueryContext(ctx context.Context, q Query, args ...interface{}) (pgx.Rows, error)
	QueryRowContext(ctx context.Context, q Query, args ...interface{}) pgx.Row
}

// Pinger интерфейс для пинга базы данных
type Pinger interface {
	Ping(ctx context.Context) error
}

// DB интерфейс базы данных
type DB interface {
	SQLExecer
	Pinger
	Close()
}
