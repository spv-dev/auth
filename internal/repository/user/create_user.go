package user

import (
	"context"
	"log"

	sq "github.com/Masterminds/squirrel"

	"github.com/spv-dev/auth/internal/client/db"
	model "github.com/spv-dev/auth/internal/model"
)

// CreateUser создаёт нового пользователя в БД
func (r *repo) CreateUser(ctx context.Context, info *model.UserInfo, password string) (int64, error) {
	builder := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns(nameColumn, emailColumn, roleColumn, passwordColumn).
		Values(info.Name, info.Email, info.Role, password).
		Suffix("returning id")

	query, args, err := builder.ToSql()
	if err != nil {
		return 0, err
	}

	q := db.Query{
		QueryRaw: query,
		Name:     "user_repository.Create",
	}

	var userID int64
	if err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&userID); err != nil {
		return 0, err
	}

	log.Printf("inserted new user with id = %v", userID)

	return userID, nil
}
