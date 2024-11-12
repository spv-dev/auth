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
	"github.com/spv-dev/auth/internal/client/kafka"
	kafkaMocks "github.com/spv-dev/auth/internal/client/kafka/mocks"
	"github.com/spv-dev/auth/internal/repository"
	repoMocks "github.com/spv-dev/auth/internal/repository/mocks"
	"github.com/spv-dev/auth/internal/service/user"
)

func TestDeleteUser(t *testing.T) {
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

		id = gofakeit.Int64()

		req     = id
		repoErr = fmt.Errorf("repo error")
	)

	tests := []struct {
		name               string
		args               args
		err                error
		userRepositoryMock userRepositoryMockFunc
		dbMockFunc         txManagerMockFunc
		userCacheMock      userCacheMockFunc
		producerMock       producerMockFunc
	}{
		{
			name: "Success Delete User",
			args: args{
				ctx: ctx,
				req: req,
			},
			err: nil,
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repoMocks.NewUserRepositoryMock(mc)
				mock.DeleteUserMock.Expect(ctx, id).Return(nil)
				return mock
			},
			dbMockFunc: func(_ *minimock.Controller) db.TxManager {
				return dbMock.NewTxManagerMock(t)
			},
			userCacheMock: func(_ *minimock.Controller) repository.UserCache {
				return repoMocks.NewUserCacheMock(t)
			},
			producerMock: func(_ *minimock.Controller) kafka.Producer {
				return kafkaMocks.NewProducerMock(t)
			},
		},
		{
			name: "Error Delete User",
			args: args{
				ctx: ctx,
				req: req,
			},
			err: repoErr,
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repoMocks.NewUserRepositoryMock(mc)
				mock.DeleteUserMock.Expect(ctx, id).Return(repoErr)
				return mock
			},
			dbMockFunc: func(_ *minimock.Controller) db.TxManager {
				return dbMock.NewTxManagerMock(t)
			},
			userCacheMock: func(_ *minimock.Controller) repository.UserCache {
				return repoMocks.NewUserCacheMock(t)
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
			userRepoMock := tt.userRepositoryMock(mc)
			txManagerMock := tt.dbMockFunc(mc)
			userCacheMock := tt.userCacheMock(mc)
			producerMock := tt.producerMock(mc)
			service := user.NewService(userRepoMock, txManagerMock, userCacheMock, producerMock)

			err := service.DeleteUser(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
		})
	}
}
