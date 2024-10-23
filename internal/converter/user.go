package converter

import (
	"database/sql"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/spv-dev/auth/internal/model"
	desc "github.com/spv-dev/auth/pkg/user_v1"
)

// ToUserFromService конвертер User из сервисного слоя в слой API
func ToUserFromService(user *model.User) *desc.User {
	var updatedAt *timestamppb.Timestamp
	if user.UpdatedAt.Valid {
		updatedAt = timestamppb.New(user.UpdatedAt.Time)
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
	return &model.UserInfo{
		Name:  info.Name,
		Email: info.Email,
		Role:  int32(info.Role),
	}
}

// ToUpdateUserInfoFromDesc конвертер UpdateUserInfo из API слоя в сервисный слой
func ToUpdateUserInfoFromDesc(info *desc.UpdateUserInfo) *model.UpdateUserInfo {
	name := sql.NullString{}
	if info.Name == nil {
		name.Valid = false
	} else {
		name = sql.NullString{
			Valid:  true,
			String: info.Name.Value,
		}
	}
	role := sql.NullInt32{}
	if info.Role == 0 {
		role.Valid = false
	} else {
		role = sql.NullInt32{
			Valid: true,
			Int32: int32(info.Role),
		}
	}
	return &model.UpdateUserInfo{
		Name: name,
		Role: role,
	}
}
