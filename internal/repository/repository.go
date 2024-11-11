package repository

import (
	"context"

	model "github.com/spv-dev/auth/internal/model"
)

// UserRepository интерфейс для взаимодействия с БД
type UserRepository interface {
	CreateUser(ctx context.Context, info *model.UserInfo, password string) (int64, error)
	GetUser(ctx context.Context, id int64) (model.User, error)
	UpdateUser(ctx context.Context, id int64, user *model.UpdateUserInfo) error
	DeleteUser(ctx context.Context, id int64) error
}

// UserCache интерфейс для кэширования пользователя
type UserCache interface {
	AddUser(ctx context.Context, id int64, user *model.User) error
	GetUser(ctx context.Context, id int64) (model.User, error)
}
