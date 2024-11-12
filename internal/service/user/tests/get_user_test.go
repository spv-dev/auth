package tests

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/spv-dev/platform_common/pkg/db"
	"github.com/stretchr/testify/require"

	dbMock "github.com/spv-dev/auth/internal/client/db/mocks"
	"github.com/spv-dev/auth/internal/client/kafka"
	kafkaMocks "github.com/spv-dev/auth/internal/client/kafka/mocks"
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
	type userCacheMockFunc func(mc *minimock.Controller) repository.UserCache
	type producerMockFunc func(mc *minimock.Controller) kafka.Producer

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
		role  = constants.RolesUSER
		dt    = time.Now()

		repoErr         = fmt.Errorf("repo error")
		addUserCacheErr = fmt.Errorf("add user cache error")

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
		userCacheMock      userCacheMockFunc
		producerMock       producerMockFunc
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
			dbMockFunc: func(_ *minimock.Controller) db.TxManager {
				return dbMock.NewTxManagerMock(t)
			},
			userCacheMock: func(_ *minimock.Controller) repository.UserCache {
				mock := repoMocks.NewUserCacheMock(t)
				mock.GetUserMock.Expect(ctx, id).Return(model.User{}, errors.New("No cache"))
				mock.AddUserMock.Expect(ctx, id, &u).Return(nil)
				return mock
			},
			producerMock: func(_ *minimock.Controller) kafka.Producer {
				return kafkaMocks.NewProducerMock(t)
			},
		},

		{
			name: "Error User Cache Add user",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: model.User{},
			err:  addUserCacheErr,
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repoMocks.NewUserRepositoryMock(mc)
				mock.GetUserMock.Expect(ctx, id).Return(u, nil)
				return mock
			},
			dbMockFunc: func(_ *minimock.Controller) db.TxManager {
				return dbMock.NewTxManagerMock(t)
			},
			userCacheMock: func(_ *minimock.Controller) repository.UserCache {
				mock := repoMocks.NewUserCacheMock(t)
				mock.GetUserMock.Expect(ctx, id).Return(model.User{}, errors.New("No cache"))
				mock.AddUserMock.Expect(ctx, id, &u).Return(addUserCacheErr)
				return mock
			},
			producerMock: func(_ *minimock.Controller) kafka.Producer {
				return kafkaMocks.NewProducerMock(t)
			},
		},

		{
			name: "Error User Cache Get user No error",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: u,
			err:  nil,
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				return repoMocks.NewUserRepositoryMock(mc)
			},
			dbMockFunc: func(_ *minimock.Controller) db.TxManager {
				return dbMock.NewTxManagerMock(t)
			},
			userCacheMock: func(_ *minimock.Controller) repository.UserCache {
				mock := repoMocks.NewUserCacheMock(t)
				mock.GetUserMock.Expect(ctx, id).Return(u, nil)

				return mock
			},
			producerMock: func(_ *minimock.Controller) kafka.Producer {
				return kafkaMocks.NewProducerMock(t)
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
			dbMockFunc: func(_ *minimock.Controller) db.TxManager {
				return dbMock.NewTxManagerMock(t)
			},
			userCacheMock: func(_ *minimock.Controller) repository.UserCache {
				mock := repoMocks.NewUserCacheMock(t)
				mock.GetUserMock.Expect(ctx, id).Return(model.User{}, errors.New("No cache"))
				return mock
			},
			producerMock: func(_ *minimock.Controller) kafka.Producer {
				return kafkaMocks.NewProducerMock(t)
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			userRepositoryMock := tt.userRepositoryMock(mc)
			txManager := tt.dbMockFunc(mc)
			cache := tt.userCacheMock(mc)
			producerMock := tt.producerMock(mc)
			service := user.NewService(userRepositoryMock, txManager, cache, producerMock)

			res, err := service.GetUser(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, res)
		})
	}
}
