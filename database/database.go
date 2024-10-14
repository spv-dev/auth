package database

import (
	"context"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	desc "github.com/spv-dev/auth/pkg/user_v1"
)

const (
	dbDSN = "host=localhost port=54321 dbname=auth user=auth-user password=authpwd sslmode=disable"
)

var pool *pgxpool.Pool

// InitDB инициализирует базу данных
func InitDB() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	p, err := pgxpool.Connect(ctx, dbDSN)
	if err != nil {
		log.Fatalf("failed to connect db: %v", err)
	}

	pool = p
}

// CloseDB закрывает базу данных
func CloseDB() {
	pool.Close()
}

// CreateUserDB создаёт пользователя в базе данных
func CreateUserDB(ctx context.Context, req *desc.CreateUserRequest) (*desc.CreateUserResponse, error) {
	builder := sq.Insert("users").
		PlaceholderFormat(sq.Dollar).
		Columns("name", "email", "role", "password").
		Values(req.Info.Name, req.Info.Email, req.Info.Role, req.Password).
		Suffix("returning id")

	query, args, err := builder.ToSql()
	if err != nil {
		log.Fatalf("failed to build query: %v", err)
	}

	var userID int64
	if err = pool.QueryRow(ctx, query, args...).Scan(&userID); err != nil {
		log.Fatalf("failed to insert user: %v", err)
	}

	log.Printf("inserted new user with id = %v", userID)

	return &desc.CreateUserResponse{
		Id: userID,
	}, nil
}

// GetUserDB получает пользователя из базы данных
func GetUserDB(ctx context.Context, req *desc.GetUserRequest) (*desc.GetUserResponse, error) {
	builder := sq.Select("id", "name", "email", "role", "created_at").
		PlaceholderFormat(sq.Dollar).
		From("users").
		Where(sq.Eq{"id": req.Id})

	query, args, err := builder.ToSql()
	if err != nil {
		log.Fatalf("failed to build query: %v", err)
	}

	u := desc.User{Info: &desc.UserInfo{}}
	var createdAt time.Time
	if err := pool.QueryRow(ctx, query, args...).
		Scan(&u.Id, &u.Info.Name, &u.Info.Email, &u.Info.Role, &createdAt); err != nil {
		log.Fatalf("failed to get user: %v", err)
	}
	u.CreatedAt = timestamppb.New(createdAt)

	log.Printf("get user info: %v", req.Id)

	return &desc.GetUserResponse{
		User: &u,
	}, nil
}

// UpdateUserDB изменяет пользователя в базе данных
func UpdateUserDB(ctx context.Context, req *desc.UpdateUserRequest) (*emptypb.Empty, error) {
	builder := sq.Update("users").
		PlaceholderFormat(sq.Dollar).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": req.Id})

	if req.Info != nil && req.Info.Name != nil && req.Info.Name.Value != "" {
		builder.Set("name", req.Info.Name.Value)
	}

	if req.Info != nil && req.Info.Role != desc.Roles_UNKNOWN {
		builder.Set("role", req.Info.Role)
	}

	query, args, err := builder.ToSql()
	if err != nil {
		log.Fatalf("failed to build query: %v", err)
	}

	res, err := pool.Exec(ctx, query, args...)
	if err != nil {
		log.Fatalf("failed to update user: %v", err)
	}

	log.Printf("updated users count: %v", res)

	return nil, nil
}

// DeleteUserDB удаляет пользователя в базе данных по идентификатору
func DeleteUserDB(ctx context.Context, req *desc.DeleteUserRequest) (*emptypb.Empty, error) {
	builder := sq.Delete("users").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": req.Id})

	query, args, err := builder.ToSql()
	if err != nil {
		log.Fatalf("failed to build query: %v", err)
	}

	res, err := pool.Exec(ctx, query, args...)
	if err != nil {
		log.Fatalf("failed to delete user: %v", err)
	}

	log.Printf("deleted users count: %v", res)

	return nil, nil
}
