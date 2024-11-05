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

func TestCreateUser(t *testing.T) {
	t.Parallel()
	type userRepositoryMockFunc func(mc *minimock.Controller) repository.UserRepository
	type txManagerMockFunc func(mc *minimock.Controller) db.TxManager
	type userCacheMockFunc func(mc *minimock.Controller) repository.UserCache

	type args struct {
		ctx      context.Context
		req      *model.UserInfo
		password string
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id    = gofakeit.Int64()
		name  = gofakeit.Name()
		email = gofakeit.Email()
		role  = constants.RolesUSER
		pass  = gofakeit.Password(true, true, true, true, false, 8)

		repoErr = fmt.Errorf("repo error")

		req = &model.UserInfo{
			Name:  name,
			Email: email,
			Role:  role,
		}

		info = &model.UserInfo{
			Name:  name,
			Email: email,
			Role:  role,
		}

		res = id
	)

	tests := []struct {
		name               string
		args               args
		want               int64
		err                error
		userRepositoryMock userRepositoryMockFunc
		dbMockFunc         txManagerMockFunc
		userCacheMock      userCacheMockFunc
	}{
		{
			name: "Success Create User",
			args: args{
				ctx:      ctx,
				req:      req,
				password: pass,
			},
			want: res,
			err:  nil,
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repoMocks.NewUserRepositoryMock(mc)
				mock.CreateUserMock.Expect(ctx, info, pass).Return(id, nil)
				return mock
			},
			dbMockFunc: func(_ *minimock.Controller) db.TxManager {
				mock := dbMock.NewTxManagerMock(t)
				mock.ReadCommitedMock.Set(func(ctx context.Context, handler db.Handler) error {
					return handler(ctx)
				})

				return mock
			},
			userCacheMock: func(_ *minimock.Controller) repository.UserCache {
				return repoMocks.NewUserCacheMock(t)
			},
		},
		{
			name: "Error Empty Info",
			args: args{
				ctx:      ctx,
				req:      nil,
				password: pass,
			},
			want: 0,
			err:  fmt.Errorf("Пустые данные при создании пользователя"),
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
			name: "Error Empty Password",
			args: args{
				ctx:      ctx,
				req:      req,
				password: "",
			},
			want: 0,
			err:  fmt.Errorf("Пустой пароль"),
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
			name: "Error Empty UserName",
			args: args{
				ctx: ctx,
				req: &model.UserInfo{
					Name:  "",
					Email: "aaa@aa.ru",
					Role:  constants.RolesUSER,
				},
				password: "112233",
			},
			want: 0,
			err:  fmt.Errorf("Пустое имя пользователя"),
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
			name: "Error Email",
			args: args{
				ctx: ctx,
				req: &model.UserInfo{
					Name:  "Name ",
					Email: "aaa dfs",
					Role:  constants.RolesUSER,
				},
				password: "1122233",
			},
			want: 0,
			err:  fmt.Errorf("Указан неверный Email"),
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
			name: "Error Email",
			args: args{
				ctx: ctx,
				req: &model.UserInfo{
					Name:  "Name ",
					Email: "",
					Role:  constants.RolesUSER,
				},
				password: "1122233",
			},
			want: 0,
			err:  fmt.Errorf("Пустой email пользователя"),
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
			name: "Error in repo",
			args: args{
				ctx:      ctx,
				req:      req,
				password: pass,
			},
			want: 0,
			err:  repoErr,
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repoMocks.NewUserRepositoryMock(mc)
				mock.CreateUserMock.Expect(ctx, info, pass).Return(0, repoErr)
				return mock
			},
			dbMockFunc: func(_ *minimock.Controller) db.TxManager {
				mock := dbMock.NewTxManagerMock(t)
				mock.ReadCommitedMock.Set(func(ctx context.Context, handler db.Handler) error {
					return handler(ctx)
				})

				return mock
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
			userRepoMock := tt.userRepositoryMock(mc)
			txManagerMock := tt.dbMockFunc(mc)
			userCacheMock := tt.userCacheMock(mc)
			service := user.NewService(userRepoMock, txManagerMock, userCacheMock)

			res, err := service.CreateUser(tt.args.ctx, tt.args.req, tt.args.password)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, res)
		})
	}
}
