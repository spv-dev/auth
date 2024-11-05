package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/spv-dev/platform_common/pkg/db"
	"github.com/stretchr/testify/require"

	dbMock "github.com/spv-dev/auth/internal/client/db/mocks"
	"github.com/spv-dev/auth/internal/constants"
	model "github.com/spv-dev/auth/internal/model"
	"github.com/spv-dev/auth/internal/repository"
	repoMocks "github.com/spv-dev/auth/internal/repository/mocks"
	"github.com/spv-dev/auth/internal/service/user"
)

func TestUpdateUser(t *testing.T) {
	t.Parallel()
	type userRepositoryMockFunc func(mc *minimock.Controller) repository.UserRepository
	type txManagerMockFunc func(mc *minimock.Controller) db.TxManager
	type userCacheMockFunc func(mc *minimock.Controller) repository.UserCache

	type args struct {
		ctx context.Context
		req *model.UpdateUserInfo
		id  int64
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id          = gofakeit.Int64()
		name        = gofakeit.Name()
		role        = constants.RolesUSER
		emptyString = ""
		repoErr     = fmt.Errorf("repo error")

		req = &model.UpdateUserInfo{
			Name: &name,
			Role: &role,
		}

		info = &model.UpdateUserInfo{
			Name: &name,
			Role: &role,
		}
	)

	tests := []struct {
		name               string
		args               args
		err                error
		userRepositoryMock userRepositoryMockFunc
		dbMockFunc         txManagerMockFunc
		userCacheMock      userCacheMockFunc
	}{
		{
			name: "Success Update User",
			args: args{
				ctx: ctx,
				req: req,
				id:  id,
			},
			err: nil,
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repoMocks.NewUserRepositoryMock(mc)
				mock.UpdateUserMock.Expect(ctx, id, info).Return(nil)
				return mock
			},
			dbMockFunc: func(_ *minimock.Controller) db.TxManager {
				return dbMock.NewTxManagerMock(t)
			},
			userCacheMock: func(_ *minimock.Controller) repository.UserCache {
				return repoMocks.NewUserCacheMock(t)
			},
		},
		{
			name: "Error Empty data",
			args: args{
				ctx: ctx,
				req: nil,
				id:  id,
			},
			err: fmt.Errorf("Пустые данные при изменении пользователя"),
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				return repoMocks.NewUserRepositoryMock(mc)
			},
			dbMockFunc: func(_ *minimock.Controller) db.TxManager {
				return dbMock.NewTxManagerMock(t)
			},
			userCacheMock: func(_ *minimock.Controller) repository.UserCache {
				return repoMocks.NewUserCacheMock(t)
			},
		},
		{
			name: "Error Empty Name",
			args: args{
				ctx: ctx,
				req: &model.UpdateUserInfo{
					Name: &emptyString,
					Role: &role,
				},
				id: id,
			},
			err: fmt.Errorf("Пустое имя пользователя"),
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				return repoMocks.NewUserRepositoryMock(mc)
			},
			dbMockFunc: func(_ *minimock.Controller) db.TxManager {
				return dbMock.NewTxManagerMock(t)
			},
			userCacheMock: func(_ *minimock.Controller) repository.UserCache {
				return repoMocks.NewUserCacheMock(t)
			},
		},
		{
			name: "Error Update User",
			args: args{
				ctx: ctx,
				req: req,
				id:  id,
			},
			err: repoErr,
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repoMocks.NewUserRepositoryMock(mc)
				mock.UpdateUserMock.Expect(ctx, id, info).Return(repoErr)
				return mock
			},
			dbMockFunc: func(_ *minimock.Controller) db.TxManager {
				return dbMock.NewTxManagerMock(t)
			},
			userCacheMock: func(_ *minimock.Controller) repository.UserCache {
				return repoMocks.NewUserCacheMock(t)
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			userRepositoryMock := tt.userRepositoryMock(mc)
			txManager := tt.dbMockFunc(mc)
			userCache := tt.userCacheMock(mc)
			service := user.NewService(userRepositoryMock, txManager, userCache)

			err := service.UpdateUser(tt.args.ctx, tt.args.id, tt.args.req)
			require.Equal(t, tt.err, err)
		})
	}
}
