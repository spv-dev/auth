package user

import (
	"github.com/spv-dev/auth/internal/repository"
	"github.com/spv-dev/auth/internal/service"
)

type serv struct {
	userRepository repository.UserRepository
}

// NewService создаёт новый сервис
func NewService(userRepository repository.UserRepository) service.UserService {
	return &serv{
		userRepository: userRepository,
	}
}
