package cache

import (
	"context"
	"errors"
	"strconv"

	"github.com/gomodule/redigo/redis"
	client "github.com/spv-dev/auth/internal/client/cache"
	"github.com/spv-dev/auth/internal/model"
	"github.com/spv-dev/auth/internal/repository"
	conv "github.com/spv-dev/auth/internal/repository/cache/converter"
	cacheModel "github.com/spv-dev/auth/internal/repository/cache/model"
)

var _ repository.UserCache = (*cache)(nil)

type cache struct {
	redisClient client.RedisClient
}

func NewCache(client client.RedisClient) repository.UserCache {
	return &cache{
		redisClient: client,
	}
}

func (c *cache) AddUser(ctx context.Context, id int64, user *model.User) error {
	if user == nil {
		return errors.New("нет задана информация о пользователе")
	}

	redisUser := conv.ToCacheFromModel(user)

	err := c.redisClient.HashSet(ctx, strconv.FormatInt(id, 10), redisUser)
	if err != nil {
		return err
	}

	return nil
}

func (c *cache) GetUser(ctx context.Context, id int64) (model.User, error) {
	values, err := c.redisClient.HGetAll(ctx, strconv.FormatInt(id, 10))
	if err != nil {
		return model.User{}, err
	}

	if len(values) == 0 {
		return model.User{}, errors.New("пользователь не найден")
	}

	var user cacheModel.UserCache
	err = redis.ScanStruct(values, &user)
	if err != nil {
		return model.User{}, err
	}

	return conv.ToModelFromCache(&user), nil
}
