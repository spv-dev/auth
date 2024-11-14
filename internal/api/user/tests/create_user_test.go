package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"

	"github.com/spv-dev/auth/internal/api/user"
	"github.com/spv-dev/auth/internal/constants"
	"github.com/spv-dev/auth/internal/converter"
	model "github.com/spv-dev/auth/internal/model"
	serviceMocks "github.com/spv-dev/auth/internal/service/mocks"
	desc "github.com/spv-dev/auth/pkg/user_v1"
)

func TestCreateUser(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	id := gofakeit.Int64()
	name := gofakeit.Name()
	email := gofakeit.Email()
	role := constants.RolesUSER
	pass := gofakeit.Password(true, true, true, false, false, 8)

	req := &desc.CreateUserRequest{
		Info: &desc.UserInfo{
			Name:  name,
			Email: email,
			Role:  converter.ConvertRoleFromModel(role),
		},
		Password: pass,
	}

	info := &model.UserInfo{
		Name:  name,
		Email: email,
		Role:  role,
	}

	mc := minimock.NewController(t)

	service := serviceMocks.NewUserServiceMock(mc)

	api := user.NewServer(service)

	t.Run("create user success", func(t *testing.T) {
		t.Parallel()

		service.CreateUserMock.Expect(ctx, info, pass).Return(id, nil)

		resp, err := api.CreateUser(ctx, req)

		assert.NoError(t, err)
		assert.Equal(t, resp.Id, id)
	})

	serviceErr := errors.New("service error")
	t.Run("create user error", func(t *testing.T) {
		t.Parallel()

		service.CreateUserMock.Expect(ctx, info, pass).Return(0, serviceErr)

		_, err := api.CreateUser(ctx, req)

		assert.Equal(t, err, serviceErr)
	})
}
