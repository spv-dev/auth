package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/spv-dev/auth/internal/api/user"
	"github.com/spv-dev/auth/internal/constants"
	"github.com/spv-dev/auth/internal/converter"
	model "github.com/spv-dev/auth/internal/model"
	serviceMocks "github.com/spv-dev/auth/internal/service/mocks"
	desc "github.com/spv-dev/auth/pkg/user_v1"
)

func TestUpdateUser(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	id := gofakeit.Int64()
	name := gofakeit.Name()
	role := constants.RolesUSER

	req := &desc.UpdateUserRequest{
		Id: id,
		Info: &desc.UpdateUserInfo{
			Name: wrapperspb.String(name),
			Role: converter.ConvertRoleFromModel(role),
		},
	}

	userInfo := &model.UpdateUserInfo{
		Name: &name,
		Role: &role,
	}

	mc := minimock.NewController(t)

	service := serviceMocks.NewUserServiceMock(mc)

	api := user.NewServer(service)

	t.Run("update user success", func(t *testing.T) {
		t.Parallel()

		service.UpdateUserMock.Expect(ctx, id, userInfo).Return(nil)

		_, err := api.UpdateUser(ctx, req)

		assert.NoError(t, err)
	})

	errorService := errors.New("service error")
	t.Run("update user error", func(t *testing.T) {
		t.Parallel()

		service.UpdateUserMock.Expect(ctx, id, userInfo).Return(errorService)

		_, err := api.UpdateUser(ctx, req)

		assert.Equal(t, err, errorService)
	})
}
