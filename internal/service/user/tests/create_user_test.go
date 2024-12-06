package tests

import (
	"context"
	"errors"
	"strconv"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"

	dbMock "github.com/spv-dev/auth/internal/client/db/mocks"
	kafkaMocks "github.com/spv-dev/auth/internal/client/kafka/mocks"
	"github.com/spv-dev/auth/internal/constants"
	model "github.com/spv-dev/auth/internal/model"
	repoMocks "github.com/spv-dev/auth/internal/repository/mocks"
	"github.com/spv-dev/auth/internal/service/user"
	serviceerror "github.com/spv-dev/auth/internal/service_error"
	"github.com/spv-dev/platform_common/pkg/db"
)

func TestCreateUser(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	id := gofakeit.Int64()
	name := gofakeit.Name()
	email := gofakeit.Email()
	role := constants.RolesUSER
	pass := gofakeit.Password(true, true, true, false, false, 8)

	userInfo := &model.UserInfo{
		Name:  name,
		Email: email,
		Role:  role,
	}

	userInfoEmptyName := &model.UserInfo{
		Name:  "",
		Email: email,
		Role:  role,
	}

	userInfoErrorEmail := &model.UserInfo{
		Name:  name,
		Email: "aaa aaa",
		Role:  role,
	}

	userInfoEmptyEmail := &model.UserInfo{
		Name:  name,
		Email: "",
		Role:  role,
	}

	errorTest := errors.New("test error")

	mc := minimock.NewController(t)

	repo := repoMocks.NewUserRepositoryMock(mc)
	trans := dbMock.NewTxManagerMock(mc)
	cache := repoMocks.NewUserCacheMock(mc)
	kafka := kafkaMocks.NewProducerMock(mc)

	service := user.NewService(repo, trans, cache, kafka)

	t.Run("create user success", func(t *testing.T) {
		t.Parallel()

		repo.CreateUserMock.Expect(ctx, userInfo, pass).Return(id, nil)
		trans.ReadCommitedMock.Set(func(ctx context.Context, handler db.Handler) error {
			return handler(ctx)
		})
		kafka.SendMock.Expect("topic_name", strconv.FormatInt(id, 10)).Return(nil)

		testID, err := service.CreateUser(ctx, userInfo, pass)

		assert.NoError(t, err)
		assert.Equal(t, testID, id)
	})

	t.Run("create user error", func(t *testing.T) {
		t.Parallel()

		repo.CreateUserMock.Expect(ctx, userInfo, pass).Return(0, errorTest)
		trans.ReadCommitedMock.Set(func(ctx context.Context, handler db.Handler) error {
			return handler(ctx)
		})

		_, err := service.CreateUser(ctx, userInfo, pass)

		assert.Equal(t, err, errorTest)
	})

	errEmptyUserInfo := errors.New(serviceerror.EmptyDataWhenCreateUser)
	t.Run("empty user info error", func(t *testing.T) {
		t.Parallel()

		_, err := service.CreateUser(ctx, nil, pass)

		assert.Equal(t, err, errEmptyUserInfo)
	})

	errEmptyPassword := errors.New("Пустой пароль")
	t.Run("empty password error", func(t *testing.T) {
		t.Parallel()

		_, err := service.CreateUser(ctx, userInfo, "")

		assert.Equal(t, err, errEmptyPassword)
	})

	errEmptyUserName := errors.New("Пустое имя пользователя")
	t.Run("empty user name error", func(t *testing.T) {
		t.Parallel()

		_, err := service.CreateUser(ctx, userInfoEmptyName, pass)

		assert.Equal(t, err, errEmptyUserName)
	})

	errUserEmail := errors.New("Указан неверный Email")
	t.Run("user email error", func(t *testing.T) {
		t.Parallel()

		_, err := service.CreateUser(ctx, userInfoErrorEmail, pass)

		assert.Equal(t, err, errUserEmail)
	})

	errEmptyUserEmail := errors.New("Пустой email пользователя")
	t.Run("user email error", func(t *testing.T) {
		t.Parallel()

		_, err := service.CreateUser(ctx, userInfoEmptyEmail, pass)

		assert.Equal(t, err, errEmptyUserEmail)
	})

	errorKafkaTest := errors.New("test kafka error")
	t.Run("kafka create user error", func(t *testing.T) {
		t.Parallel()

		repo.CreateUserMock.Expect(ctx, userInfo, pass).Return(id, nil)
		trans.ReadCommitedMock.Set(func(ctx context.Context, handler db.Handler) error {
			return handler(ctx)
		})
		kafka.SendMock.Expect("topic_name", strconv.FormatInt(id, 10)).Return(errorKafkaTest)

		_, err := service.CreateUser(ctx, userInfo, pass)

		assert.Equal(t, err, errorKafkaTest)
	})
}
