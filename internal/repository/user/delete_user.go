package user

import (
	"context"
	"fmt"
	"log"

	sq "github.com/Masterminds/squirrel"

	"github.com/spv-dev/auth/internal/client/db"
)

// DeleteUser удаляет пользователя в БД по идентификатору
func (r *repo) DeleteUser(ctx context.Context, id int64) error {
	builder := sq.Delete(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{idColumn: id})

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
		return fmt.Errorf("failed to delete user: %v", err)
	}

	log.Printf("deleted users count: %v", res)

	return nil
}
