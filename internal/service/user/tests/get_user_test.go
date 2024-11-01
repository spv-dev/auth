package tests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/spv-dev/platform_common/pkg/db"
	"github.com/stretchr/testify/require"

	dbMock "github.com/spv-dev/auth/internal/client/db/mocks"
	"github.com/spv-dev/auth/internal/constants"
	"github.com/spv-dev/auth/internal/model"
	"github.com/spv-dev/auth/internal/repository"
	repoMocks "github.com/spv-dev/auth/internal/repository/mocks"
	"github.com/spv-dev/auth/internal/service/user"
)

func TestGetUser(t *testing.T) {
	t.Parallel()
	type userRepositoryMockFunc func(mc *minimock.Controller) repository.UserRepository
	type txManagerMockFunc func(mc *minimock.Controller) db.TxManager

	type args struct {
		ctx context.Context
		req int64
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id    = gofakeit.Int64()
		name  = gofakeit.Name()
		email = gofakeit.Email()
		role  = constants.Roles_USER
		dt    = time.Now()

		repoErr = fmt.Errorf("repo error")

		req = id

		u = model.User{
			ID:        id,
			CreatedAt: dt,
			UpdatedAt: &dt,
			Info: model.UserInfo{
				Name:  name,
				Email: email,
				Role:  role,
			},
		}

		res = model.User{

			ID:        id,
			CreatedAt: dt,
			UpdatedAt: &dt,
			Info: model.UserInfo{
				Name:  name,
				Email: email,
				Role:  role,
			},
		}
	)

	tests := []struct {
		name               string
		args               args
		want               model.User
		err                error
		userRepositoryMock userRepositoryMockFunc
		dbMockFunc         txManagerMockFunc
	}{
		{
			name: "Success Get User",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: res,
			err:  nil,
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repoMocks.NewUserRepositoryMock(mc)
				mock.GetUserMock.Expect(ctx, id).Return(u, nil)
				return mock
			},
			dbMockFunc: func(mc *minimock.Controller) db.TxManager {
				return dbMock.NewTxManagerMock(t)
			},
		},

		{
			name: "Error Get User",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: model.User{},
			err:  repoErr,
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repoMocks.NewUserRepositoryMock(mc)
				mock.GetUserMock.Expect(ctx, id).Return(model.User{}, repoErr)
				return mock
			},
			dbMockFunc: func(mc *minimock.Controller) db.TxManager {
				return dbMock.NewTxManagerMock(t)
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			userRepositoryMock := tt.userRepositoryMock(mc)
			txManager := tt.dbMockFunc(mc)
			service := user.NewService(userRepositoryMock, txManager)

			res, err := service.GetUser(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, res)
		})
	}
}
