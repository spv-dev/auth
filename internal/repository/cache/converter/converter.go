package converter

import (
	"strconv"
	"time"

	"github.com/spv-dev/auth/internal/constants"
	model "github.com/spv-dev/auth/internal/model"
	cacheModel "github.com/spv-dev/auth/internal/repository/cache/model"
)

func ToCacheFromModel(user *model.User) cacheModel.UserCache {
	if user == nil {
		return cacheModel.UserCache{}
	}
	var updatedAt int64
	if user.UpdatedAt != nil {
		updatedAt = user.UpdatedAt.Unix()
	}

	return cacheModel.UserCache{
		ID:        strconv.FormatInt(user.ID, 10),
		Name:      user.Info.Name,
		Email:     user.Info.Email,
		Role:      int32(user.Info.Role),
		CreatedAt: user.CreatedAt.Unix(),
		UpdatedAt: &updatedAt,
	}
}

func ToModelFromCache(user *cacheModel.UserCache) model.User {
	if user == nil {
		return model.User{}
	}
	id, err := strconv.ParseInt(user.ID, 10, 64)
	if err != nil {
		panic(err)
	}
	updatedAt := time.Unix(*user.UpdatedAt, 0)
	return model.User{
		ID: id,
		Info: model.UserInfo{
			Email: user.Email,
			Role:  constants.Roles(user.Role),
			Name:  user.Name,
		},
		CreatedAt: time.Unix(user.CreatedAt, 0),
		UpdatedAt: &updatedAt,
	}
}
