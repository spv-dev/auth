package user

import (
	"context"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/protobuf/types/known/emptypb"

	model "github.com/spv-dev/auth/internal/model"
	"github.com/spv-dev/auth/internal/repository"
	"github.com/spv-dev/auth/internal/repository/user/converter"
	modelRepo "github.com/spv-dev/auth/internal/repository/user/model"
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
	db *pgxpool.Pool
}

// NewRepository получает соединение с БД
func NewRepository(db *pgxpool.Pool) repository.UserRepository {
	return &repo{db: db}
}

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

	var userID int64
	if err = r.db.QueryRow(ctx, query, args...).Scan(&userID); err != nil {
		return 0, err
	}

	log.Printf("inserted new user with id = %v", userID)

	return userID, nil
}

// GetUser получает информацию о пользователе по идентификатору
func (r *repo) GetUser(ctx context.Context, id int64) (*model.User, error) {
	builder := sq.Select(idColumn, nameColumn, emailColumn, createdAtColumn).
		PlaceholderFormat(sq.Dollar).
		From(tableName).
		Where(sq.Eq{idColumn: id}).
		Limit(1)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	var user modelRepo.User
	if err := r.db.QueryRow(ctx, query, args...).
		Scan(&user.ID, &user.Info.Name, &user.Info.Email, &user.CreatedAt); err != nil {
		return nil, err
	}

	log.Printf("get user info: %v", id)

	return converter.ToUserFromRepo(&user), nil
}

// UpdateUser изменяет пользователя в БД
func (r *repo) UpdateUser(ctx context.Context, id int64, info *model.UpdateUserInfo) (*emptypb.Empty, error) {
	builder := sq.Update(tableName).
		PlaceholderFormat(sq.Dollar).
		Set(updatedAtColumn, time.Now()).
		Where(sq.Eq{idColumn: id})

	log.Printf("id: %+v , info: %+v", id, info)

	if info != nil && info.Name.Valid {
		builder = builder.Set(nameColumn, info.Name.String)
	}

	if info != nil && info.Role.Valid {
		builder = builder.Set(roleColumn, info.Role.Int32)
	}

	query, args, err := builder.ToSql()
	if err != nil {
		log.Fatalf("failed to build query: %v", err)
	}

	log.Printf("query = %s", query)
	res, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		log.Fatalf("failed to update user: %v", err)
	}

	log.Printf("updated users count: %v", res)

	return nil, nil
}

// DeleteUser удаляет пользователя в БД по идентификатору
func (r *repo) DeleteUser(ctx context.Context, id int64) (*emptypb.Empty, error) {
	builder := sq.Delete(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{idColumn: id})

	query, args, err := builder.ToSql()
	if err != nil {
		log.Fatalf("failed to build query: %v", err)
	}

	res, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		log.Fatalf("failed to delete user: %v", err)
	}

	log.Printf("deleted users count: %v", res)

	return nil, nil
}
