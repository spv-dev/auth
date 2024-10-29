package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/spv-dev/auth/internal/api/user"
	model "github.com/spv-dev/auth/internal/model"
	"github.com/spv-dev/auth/internal/service"
	serviceMocks "github.com/spv-dev/auth/internal/service/mocks"
	desc "github.com/spv-dev/auth/pkg/user_v1"
)

func TestUpdateUser(t *testing.T) {
	type userServiceMockFunc func(mc *minimock.Controller) service.UserService

	type args struct {
		ctx context.Context
		req *desc.UpdateUserRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id   = gofakeit.Int64()
		name = gofakeit.Name()
		//role   = desc.Roles_USER
		//roleId = desc.Roles_value["USER"]

		serviceErr = fmt.Errorf("service error")

		req = &desc.UpdateUserRequest{
			Id: id,
			Info: &desc.UpdateUserInfo{
				Name: wrapperspb.String(name),
				//Role: role,
			},
		}

		info = &model.UpdateUserInfo{
			Name: &name,
			//Role: &roleId,
		}
	)

	defer t.Cleanup(mc.Finish)

	tests := []struct {
		name            string
		args            args
		want            *emptypb.Empty
		err             error
		userServiceMock userServiceMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: nil,
			err:  nil,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				mock.UpdateUserMock.Expect(ctx, id, info).Return(nil)
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
				mock.UpdateUserMock.Expect(ctx, id, info).Return(serviceErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			userServiceMock := tt.userServiceMock(mc)
			api := user.NewServer(userServiceMock)

			res, err := api.UpdateUser(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, res)
		})
	}
}
