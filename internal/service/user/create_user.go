package user

import (
	"context"

	"github.com/spv-dev/auth/internal/model"
)

// CreateUser проверяет пользователя и отправляет на создание в слой БД
func (s *serv) CreateUser(ctx context.Context, info *model.UserInfo, password string) (int64, error) {
	var id int64
	err := s.txManager.ReadCommited(ctx, func(ctx context.Context) error {
		var errTx error
		id, errTx = s.userRepository.CreateUser(ctx, info, password)
		if errTx != nil {
			return errTx
		}

		_, errTx = s.userRepository.GetUser(ctx, id)
		if errTx != nil {
			return errTx
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	return id, nil
}
