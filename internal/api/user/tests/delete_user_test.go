package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"

	"github.com/spv-dev/auth/internal/api/user"
	serviceMocks "github.com/spv-dev/auth/internal/service/mocks"
	desc "github.com/spv-dev/auth/pkg/user_v1"
)

func TestDeleteUser(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	id := gofakeit.Int64()

	req := &desc.DeleteUserRequest{
		Id: id,
	}

	mc := minimock.NewController(t)

	service := serviceMocks.NewUserServiceMock(mc)

	api := user.NewServer(service)

	t.Run("delete user success", func(t *testing.T) {
		t.Parallel()

		service.DeleteUserMock.Expect(ctx, id).Return(nil)

		_, err := api.DeleteUser(ctx, req)

		assert.NoError(t, err)
	})

	errorService := errors.New("service error")
	t.Run("delete user error", func(t *testing.T) {
		t.Parallel()

		service.DeleteUserMock.Expect(ctx, id).Return(errorService)

		_, err := api.DeleteUser(ctx, req)

		assert.Equal(t, err, errorService)
	})
}
