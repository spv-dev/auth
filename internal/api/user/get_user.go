package user

import (
	"context"
	"log"

	"github.com/spv-dev/auth/internal/converter"
	desc "github.com/spv-dev/auth/pkg/user_v1"
)

// GetUser получает информацию о пользователе по идентификатору
func (s *Server) GetUser(ctx context.Context, req *desc.GetUserRequest) (*desc.GetUserResponse, error) {
	userObj, err := s.userService.GetUser(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	log.Printf("get user by id: %d", req.GetId())

	return &desc.GetUserResponse{
		User: converter.ToUserFromService(&userObj),
	}, nil
}
