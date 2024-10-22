package user

import (
	"context"

	"github.com/spv-dev/auth/internal/model"
)

// CreateUser проверяет пользователя и отправляет на создание в слой БД
func (s *serv) CreateUser(ctx context.Context, info *model.UserInfo, password string) (int64, error) {
	id, err := s.userRepository.CreateUser(ctx, info, password)
	if err != nil {
		return 0, err
	}

	return id, nil
}
