package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	"github.com/spv-dev/auth/internal/api/user"
	"github.com/spv-dev/auth/internal/converter"
	model "github.com/spv-dev/auth/internal/model"
	"github.com/spv-dev/auth/internal/service"
	serviceMocks "github.com/spv-dev/auth/internal/service/mocks"
	desc "github.com/spv-dev/auth/pkg/user_v1"
)

func TestCreateUser(t *testing.T) {
	t.Parallel()
	type userServiceMockFunc func(mc *minimock.Controller) service.UserService

	type args struct {
		ctx context.Context
		req *desc.CreateUserRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id    = gofakeit.Int64()
		name  = gofakeit.Name()
		email = gofakeit.Email()
		role  = desc.Roles_USER
		pass  = gofakeit.Password(true, true, true, true, false, 8)

		serviceErr = fmt.Errorf("service error")

		req = &desc.CreateUserRequest{
			Info: &desc.UserInfo{
				Name:  name,
				Email: email,
				Role:  role,
			},
			Password: pass,
		}

		info = &model.UserInfo{
			Name:  name,
			Email: email,
			Role:  converter.ConvertRoleFromDesc(role),
		}

		res = &desc.CreateUserResponse{
			Id: id,
		}
	)

	tests := []struct {
		name            string
		args            args
		want            *desc.CreateUserResponse
		err             error
		userServiceMock userServiceMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: res,
			err:  nil,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				mock.CreateUserMock.Expect(ctx, info, pass).Return(id, nil)
				return mock
			},
		},
		{
			name: "error case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: nil,
			err:  serviceErr,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				mock.CreateUserMock.Expect(ctx, info, pass).Return(0, serviceErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			userServiceMock := tt.userServiceMock(mc)
			api := user.NewServer(userServiceMock)

			res, err := api.CreateUser(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, res)
		})
	}
}
