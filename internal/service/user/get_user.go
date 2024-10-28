package user

import (
	"context"

	"github.com/spv-dev/auth/internal/model"
)

// GetUser получение пользователя по идентификатору
func (s *serv) GetUser(ctx context.Context, id int64) (model.User, error) {
	user, err := s.userRepository.GetUser(ctx, id)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}
