package tests

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/spv-dev/auth/internal/api/user"
	"github.com/spv-dev/auth/internal/converter"
	"github.com/spv-dev/auth/internal/model"
	serviceMocks "github.com/spv-dev/auth/internal/service/mocks"
	desc "github.com/spv-dev/auth/pkg/user_v1"
)

func TestGetUser(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	id := gofakeit.Int64()
	name := gofakeit.Name()
	email := gofakeit.Email()
	role := desc.Roles_USER
	dt := time.Now()

	req := &desc.GetUserRequest{
		Id: id,
	}

	u := model.User{
		ID:        id,
		CreatedAt: dt,
		UpdatedAt: &dt,
		Info: model.UserInfo{
			Name:  name,
			Email: email,
			Role:  converter.ConvertRoleFromDesc(role),
		},
	}

	resp := &desc.GetUserResponse{
		User: &desc.User{
			Id:        id,
			CreatedAt: timestamppb.New(dt),
			UpdatedAt: timestamppb.New(dt),
			Info: &desc.UserInfo{
				Name:  name,
				Email: email,
				Role:  role,
			},
		},
	}

	mc := minimock.NewController(t)

	service := serviceMocks.NewUserServiceMock(mc)

	api := user.NewServer(service)

	t.Run("get user success", func(t *testing.T) {
		t.Parallel()

		service.GetUserMock.Expect(ctx, id).Return(u, nil)

		user, err := api.GetUser(ctx, req)

		assert.NoError(t, err)
		assert.Equal(t, user, resp)
	})

	errorService := errors.New("service error")
	t.Run("get user error", func(t *testing.T) {
		t.Parallel()

		service.GetUserMock.Expect(ctx, id).Return(model.User{}, errorService)

		_, err := api.GetUser(ctx, req)

		assert.Equal(t, err, errorService)
	})
}
