package repository

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	desc "github.com/spv-dev/auth/pkg/user_v1"
)

// UserRepository интерфейс для взаимодействия с БД
type UserRepository interface {
	CreateUser(ctx context.Context, info *desc.UserInfo, password string) (int64, error)
	GetUser(ctx context.Context, id int64) (*desc.User, error)
	UpdateUser(ctx context.Context, id int64, user *desc.UpdateUserInfo) (*emptypb.Empty, error)
	DeleteUser(ctx context.Context, id int64) (*emptypb.Empty, error)
}
