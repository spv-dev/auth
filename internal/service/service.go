package service

import (
	"context"

	"github.com/spv-dev/auth/internal/model"
	"google.golang.org/protobuf/types/known/emptypb"
)

// UserService сервисный слой
type UserService interface {
	CreateUser(ctx context.Context, info *model.UserInfo, password string) (int64, error)
	GetUser(ctx context.Context, id int64) (*model.User, error)
	UpdateUser(ctx context.Context, id int64, info *model.UpdateUserInfo) (*emptypb.Empty, error)
	DeleteUser(ctx context.Context, id int64) (*emptypb.Empty, error)
}
