package user

import (
	"context"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/spv-dev/platform_common/pkg/db"

	model "github.com/spv-dev/auth/internal/model"
)

// GetUser получает информацию о пользователе по идентификатору
func (r *repo) GetUser(ctx context.Context, id int64) (model.User, error) {
	builder := sq.Select(idColumn, nameColumn, emailColumn, createdAtColumn).
		From(tableName).
		Where(sq.Eq{idColumn: id}).
		Limit(1).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return model.User{}, err
	}

	q := db.Query{
		QueryRaw: query,
		Name:     "user_repository.Get",
	}

	var user model.User
	if err := r.db.DB().ScanOneContext(ctx, &user, q, args...); err != nil {
		return model.User{}, err
	}

	log.Printf("get user info: %v", id)

	return user, nil
}
