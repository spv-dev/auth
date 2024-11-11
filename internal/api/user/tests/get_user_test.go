package tests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/spv-dev/auth/internal/api/user"
	"github.com/spv-dev/auth/internal/converter"
	"github.com/spv-dev/auth/internal/model"
	"github.com/spv-dev/auth/internal/service"
	serviceMocks "github.com/spv-dev/auth/internal/service/mocks"
	desc "github.com/spv-dev/auth/pkg/user_v1"
)

func TestGetUser(t *testing.T) {
	t.Parallel()
	type userServiceMockFunc func(mc *minimock.Controller) service.UserService

	type args struct {
		ctx context.Context
		req *desc.GetUserRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id    = gofakeit.Int64()
		name  = gofakeit.Name()
		email = gofakeit.Email()
		role  = desc.Roles_USER
		dt    = time.Now()

		serviceErr = fmt.Errorf("service error")

		req = &desc.GetUserRequest{
			Id: id,
		}

		u = model.User{
			ID:        id,
			CreatedAt: dt,
			UpdatedAt: &dt,
			Info: model.UserInfo{
				Name:  name,
				Email: email,
				Role:  converter.ConvertRoleFromDesc(role),
			},
		}

		res = &desc.GetUserResponse{
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
	)

	tests := []struct {
		name            string
		args            args
		want            *desc.GetUserResponse
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
				mock.GetUserMock.Expect(ctx, id).Return(u, nil)
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
				mock.GetUserMock.Expect(ctx, id).Return(model.User{}, serviceErr)
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

			res, err := api.GetUser(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, res)
		})
	}
}
