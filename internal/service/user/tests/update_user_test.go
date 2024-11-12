package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"

	dbMock "github.com/spv-dev/auth/internal/client/db/mocks"
	kafkaMocks "github.com/spv-dev/auth/internal/client/kafka/mocks"
	"github.com/spv-dev/auth/internal/constants"
	model "github.com/spv-dev/auth/internal/model"
	repoMocks "github.com/spv-dev/auth/internal/repository/mocks"
	"github.com/spv-dev/auth/internal/service/user"
)

func TestUpdateUser(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	id := gofakeit.Int64()
	name := gofakeit.Name()
	role := constants.RolesUSER

	userInfo := &model.UpdateUserInfo{
		Name: &name,
		Role: &role,
	}

	emptyString := ""

	userInfoEmptyName := &model.UpdateUserInfo{
		Name: &emptyString,
		Role: &role,
	}

	repoError := errors.New("repo error")

	mc := minimock.NewController(t)

	repo := repoMocks.NewUserRepositoryMock(mc)
	trans := dbMock.NewTxManagerMock(mc)
	cache := repoMocks.NewUserCacheMock(mc)
	kafka := kafkaMocks.NewProducerMock(mc)

	service := user.NewService(repo, trans, cache, kafka)

	t.Run("update user success", func(t *testing.T) {
		t.Parallel()

		repo.UpdateUserMock.Expect(ctx, id, userInfo).Return(nil)

		err := service.UpdateUser(ctx, id, userInfo)
		assert.NoError(t, err)
	})

	t.Run("update user error", func(t *testing.T) {
		t.Parallel()

		repo.UpdateUserMock.Expect(ctx, id, userInfo).Return(repoError)

		err := service.UpdateUser(ctx, id, userInfo)
		assert.ErrorIs(t, err, repoError)
	})

	emptyErr := errors.New("Пустые данные при изменении пользователя")
	t.Run("update user error empty data", func(t *testing.T) {
		//t.Parallel()

		repo.UpdateUserMock.Expect(ctx, id, nil).Return(emptyErr)

		err := service.UpdateUser(ctx, id, nil)
		assert.Equal(t, err, emptyErr)
	})

	errEmptyName := errors.New("Пустое имя пользователя")
	t.Run("update user error empty name", func(t *testing.T) {
		t.Parallel()

		repo.UpdateUserMock.Expect(ctx, id, userInfoEmptyName).Return(errEmptyName)

		err := service.UpdateUser(ctx, id, userInfoEmptyName)
		assert.Equal(t, err, errEmptyName)
	})

}
