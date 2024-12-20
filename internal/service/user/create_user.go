package user

import (
	"context"
	"fmt"

	"github.com/spv-dev/auth/internal/model"
	serviceerror "github.com/spv-dev/auth/internal/service_error"
	"github.com/spv-dev/auth/internal/validator"
)

// CreateUser проверяет пользователя и отправляет на создание в слой БД
func (s *serv) CreateUser(ctx context.Context, info *model.UserInfo, password string) (int64, error) {
	if info == nil {
		return 0, fmt.Errorf(serviceerror.EmptyDataWhenCreateUser)
	}

	// проверки
	if err := validator.CheckName(info.Name); err != nil {
		return 0, err
	}

	if err := validator.CheckEmail(info.Email); err != nil {
		return 0, err
	}

	if err := validator.CheckPassword(password); err != nil {
		return 0, err
	}

	var id int64
	err := s.txManager.ReadCommited(ctx, func(ctx context.Context) error {
		var errTx error
		id, errTx = s.userRepository.CreateUser(ctx, info, password)
		if errTx != nil {
			return errTx
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	err = s.userProducer.Send("topic_name", fmt.Sprintf("%d", id))
	if err != nil {
		return 0, err
	}

	return id, nil
}
