package user

import (
	"context"
	"errors"

	"github.com/spv-dev/auth/internal/model"
	"github.com/spv-dev/auth/internal/validator"
)

// UpdateUser изменение данных о пользователе
func (s *serv) UpdateUser(ctx context.Context, id int64, info *model.UpdateUserInfo) error {
	if info == nil {
		return errors.New("Пустые данные при изменении пользователя")
	}

	// проверки
	if info.Name != nil {
		if err := validator.CheckName(*info.Name); err != nil {
			return err
		}
	}

	err := s.userRepository.UpdateUser(ctx, id, info)
	if err != nil {
		return err
	}

	return nil
}
