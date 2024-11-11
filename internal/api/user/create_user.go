package user

import (
	"context"
	"log"

	"github.com/spv-dev/auth/internal/converter"
	desc "github.com/spv-dev/auth/pkg/user_v1"
)

// CreateUser создаёт нового пользователя
func (s *Server) CreateUser(ctx context.Context, req *desc.CreateUserRequest) (*desc.CreateUserResponse, error) {
	id, err := s.userService.CreateUser(ctx, converter.ToUserInfoFromDesc(req.GetInfo()), req.Password)
	if err != nil {
		return nil, err
	}

	log.Printf("inserted user with id: %d", id)

	return &desc.CreateUserResponse{
		Id: id,
	}, nil
}
