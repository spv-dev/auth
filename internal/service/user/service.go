package user

import (
	"github.com/spv-dev/platform_common/pkg/db"

	"github.com/spv-dev/auth/internal/repository"
	"github.com/spv-dev/auth/internal/service"
)

type serv struct {
	userRepository repository.UserRepository
	userCache      repository.UserCache
	txManager      db.TxManager
}

// NewService создаёт новый сервис
func NewService(userRepository repository.UserRepository,
	txManager db.TxManager,
	userCache repository.UserCache) service.UserService {
	return &serv{
		userRepository: userRepository,
		userCache:      userCache,
		txManager:      txManager,
	}
}
