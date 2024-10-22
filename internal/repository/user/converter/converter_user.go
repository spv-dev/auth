package converter

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/spv-dev/auth/internal/repository/user/model"
	desc "github.com/spv-dev/auth/pkg/user_v1"
)

// ToUserFromRepo преобразует model.User в User
func ToUserFromRepo(user *model.User) *desc.User {
	var updatedAt *timestamppb.Timestamp
	if user.UpdatedAt.Valid {
		updatedAt = timestamppb.New(user.UpdatedAt.Time)
	}

	return &desc.User{
		Id:        user.ID,
		Info:      ToUserInfoFromRepo(user.Info),
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: updatedAt,
	}
}

// ToUserInfoFromRepo преобразует model.Info в UserInfo
func ToUserInfoFromRepo(info model.Info) *desc.UserInfo {
	return &desc.UserInfo{
		Name:  info.Name,
		Email: info.Email,
		Role:  desc.Roles(info.Role),
	}
}
