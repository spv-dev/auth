package repository

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	model "github.com/spv-dev/auth/internal/model"
)

// UserRepository интерфейс для взаимодействия с БД
type UserRepository interface {
	CreateUser(ctx context.Context, info *model.UserInfo, password string) (int64, error)
	GetUser(ctx context.Context, id int64) (*model.User, error)
	UpdateUser(ctx context.Context, id int64, user *model.UpdateUserInfo) (*emptypb.Empty, error)
	DeleteUser(ctx context.Context, id int64) (*emptypb.Empty, error)
}
