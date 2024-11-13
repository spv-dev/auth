package tests

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"

	dbMock "github.com/spv-dev/auth/internal/client/db/mocks"
	kafkaMocks "github.com/spv-dev/auth/internal/client/kafka/mocks"
	"github.com/spv-dev/auth/internal/constants"
	"github.com/spv-dev/auth/internal/model"
	repoMocks "github.com/spv-dev/auth/internal/repository/mocks"
	"github.com/spv-dev/auth/internal/service/user"
)

func TestGetUser(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	id := gofakeit.Int64()
	name := gofakeit.Name()
	email := gofakeit.Email()
	role := constants.RolesUSER
	dt := time.Now()

	errorTest := errors.New("test error")

	mc := minimock.NewController(t)

	repo := repoMocks.NewUserRepositoryMock(mc)
	trans := dbMock.NewTxManagerMock(mc)
	cache := repoMocks.NewUserCacheMock(mc)
	kafka := kafkaMocks.NewProducerMock(mc)

	service := user.NewService(repo, trans, cache, kafka)

	user := model.User{
		ID:        id,
		CreatedAt: dt,
		UpdatedAt: &dt,
		Info: model.UserInfo{
			Name:  name,
			Email: email,
			Role:  role,
		},
	}

	errCacheNotFound := errors.New("пользователь не найден")
	t.Run("get user db success", func(t *testing.T) {
		t.Parallel()

		repo.GetUserMock.Expect(ctx, id).Return(user, nil)
		cache.AddUserMock.Expect(ctx, id, &user).Return(nil)
		cache.GetUserMock.Expect(ctx, id).Return(model.User{}, errCacheNotFound)

		u, err := service.GetUser(ctx, id)

		assert.NoError(t, err)
		assert.Equal(t, u, user)
	})

	t.Run("get user db error", func(t *testing.T) {
		t.Parallel()

		repo.GetUserMock.Expect(ctx, id).Return(model.User{}, errorTest)
		cache.GetUserMock.Expect(ctx, id).Return(model.User{}, errCacheNotFound)

		_, err := service.GetUser(ctx, id)

		assert.Equal(t, err, errorTest)
	})

	t.Run("get user cache success", func(t *testing.T) {
		t.Parallel()
		cache.GetUserMock.Expect(ctx, id).Return(user, nil)

		u, err := service.GetUser(ctx, id)

		assert.NoError(t, err)
		assert.Equal(t, u, user)
	})

	errNotUserInfo := errors.New("нет задана информация о пользователе")
	t.Run("get user add cache error", func(t *testing.T) {
		t.Parallel()
		repo.GetUserMock.Expect(ctx, id).Return(user, nil)
		cache.AddUserMock.Expect(ctx, id, &user).Return(errNotUserInfo)
		cache.GetUserMock.Expect(ctx, id).Return(model.User{}, errCacheNotFound)

		_, err := service.GetUser(ctx, id)

		assert.Equal(t, err, errNotUserInfo)
	})
}
