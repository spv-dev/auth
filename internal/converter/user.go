package converter

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/spv-dev/auth/internal/model"
	desc "github.com/spv-dev/auth/pkg/user_v1"
)

// ToUserFromService конвертер User из сервисного слоя в слой API
func ToUserFromService(user *model.User) *desc.User {
	if user == nil {
		return &desc.User{}
	}
	var updatedAt *timestamppb.Timestamp
	if user.UpdatedAt != nil {
		updatedAt = timestamppb.New(*user.UpdatedAt)
	}
	return &desc.User{
		Id:   user.ID,
		Info: ToUserInfoFromService(user.Info),
		//Password:  user.Password,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: updatedAt,
	}
}

// ToUserInfoFromService конвертер UserInfo из сервисного слоя в слой API
func ToUserInfoFromService(info model.UserInfo) *desc.UserInfo {
	return &desc.UserInfo{
		Name:  info.Name,
		Email: info.Email,
		Role:  desc.Roles(info.Role),
	}
}

// ToUserInfoFromDesc конвертер UserInfo из API слоя в сервисный слой
func ToUserInfoFromDesc(info *desc.UserInfo) *model.UserInfo {
	if info == nil {
		return &model.UserInfo{}
	}
	return &model.UserInfo{
		Name:  info.Name,
		Email: info.Email,
		Role:  int32(info.Role),
	}
}

// ToUpdateUserInfoFromDesc конвертер UpdateUserInfo из API слоя в сервисный слой
func ToUpdateUserInfoFromDesc(info *desc.UpdateUserInfo) *model.UpdateUserInfo {
	var userInfo model.UpdateUserInfo
	if info == nil {
		return &userInfo
	}

	if info.Name != nil {
		userInfo.Name = &info.Name.Value
	}
	if info.Role != 0 {
		userInfo.Name = &info.Name.Value
	}

	return &userInfo
}
