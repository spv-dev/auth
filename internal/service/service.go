package service

import (
	"context"

	"github.com/spv-dev/auth/internal/model"
)

// UserService сервисный слой
type UserService interface {
	CreateUser(ctx context.Context, info *model.UserInfo, password string) (int64, error)
	GetUser(ctx context.Context, id int64) (model.User, error)
	UpdateUser(ctx context.Context, id int64, info *model.UpdateUserInfo) error
	DeleteUser(ctx context.Context, id int64) error
}
