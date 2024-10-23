package user

import (
	"github.com/spv-dev/auth/internal/client/db"
	"github.com/spv-dev/auth/internal/repository"
	"github.com/spv-dev/auth/internal/service"
)

type serv struct {
	userRepository repository.UserRepository
	txManager      db.TxManager
}

// NewService создаёт новый сервис
func NewService(userRepository repository.UserRepository,
	txManager db.TxManager) service.UserService {
	return &serv{
		userRepository: userRepository,
		txManager:      txManager,
	}
}
