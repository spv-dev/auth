package user

import (
	"github.com/spv-dev/auth/internal/client/db"
	"github.com/spv-dev/auth/internal/repository"
)

const (
	tableName = "users"

	idColumn        = "id"
	nameColumn      = "name"
	emailColumn     = "email"
	roleColumn      = "role"
	passwordColumn  = "password"
	createdAtColumn = "created_at"
	updatedAtColumn = "updated_at"
)

type repo struct {
	db db.Client
}

// NewRepository получает соединение с БД
func NewRepository(db db.Client) repository.UserRepository {
	return &repo{db: db}
}
