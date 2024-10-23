package converter

import (
	"github.com/spv-dev/auth/internal/model"
	modelRepo "github.com/spv-dev/auth/internal/repository/user/model"
)

// ToUserFromRepo преобразует model.User в User
func ToUserFromRepo(user *modelRepo.User) *model.User {

	return &model.User{
		ID:        user.ID,
		Info:      ToUserInfoFromRepo(user.Info),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

// ToUserInfoFromRepo преобразует model.Info в UserInfo
func ToUserInfoFromRepo(info modelRepo.UserInfo) model.UserInfo {
	return model.UserInfo{
		Name:  info.Name,
		Email: info.Email,
		Role:  info.Role,
	}
}
