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
	repoMocks "github.com/spv-dev/auth/internal/repository/mocks"
	"github.com/spv-dev/auth/internal/service/user"
)

func TestDeleteUser(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	id := gofakeit.Int64()

	testError := errors.New("test error")

	mc := minimock.NewController(t)

	repo := repoMocks.NewUserRepositoryMock(mc)
	trans := dbMock.NewTxManagerMock(mc)
	cache := repoMocks.NewUserCacheMock(mc)
	kafka := kafkaMocks.NewProducerMock(mc)

	service := user.NewService(repo, trans, cache, kafka)

	t.Run("delete user success", func(t *testing.T) {
		t.Parallel()

		repo.DeleteUserMock.Expect(ctx, id).Return(nil)

		err := service.DeleteUser(ctx, id)
		assert.NoError(t, err)
	})

	t.Run("delete user error", func(t *testing.T) {
		t.Parallel()

		repo.DeleteUserMock.Expect(ctx, id).Return(testError)

		err := service.DeleteUser(ctx, id)
		assert.ErrorIs(t, err, testError)
	})
}
