package user

import (
	"context"
	"fmt"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/spv-dev/auth/internal/client/db"
	model "github.com/spv-dev/auth/internal/model"
)

// UpdateUser изменяет пользователя в БД
func (r *repo) UpdateUser(ctx context.Context, id int64, info *model.UpdateUserInfo) error {
	builder := sq.Update(tableName).
		PlaceholderFormat(sq.Dollar).
		Set(updatedAtColumn, time.Now()).
		Where(sq.Eq{idColumn: id})

	if info.Name != nil {
		builder = builder.Set(nameColumn, info.Name)
	}

	if info.Role != nil {
		builder = builder.Set(roleColumn, info.Role)
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %v", err)
	}

	q := db.Query{
		QueryRaw: query,
		Name:     "user_repository.Update",
	}

	res, err := r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return fmt.Errorf("failed to update user: %v", err)
	}

	log.Printf("updated users count: %v", res)

	return nil
}
