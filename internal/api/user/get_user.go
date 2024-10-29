package user

import (
	"context"
	"log"

	"github.com/spv-dev/auth/internal/converter"
	desc "github.com/spv-dev/auth/pkg/user_v1"
)

// GetUser получает информацию о пользователе по идентификатору
func (s *Server) GetUser(ctx context.Context, req *desc.GetUserRequest) (*desc.GetUserResponse, error) {
	id := req.GetId()
	userObj, err := s.userService.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}

	log.Printf("get user by id: %d", id)

	user := converter.ToUserFromService(&userObj)

	return &desc.GetUserResponse{
		User: &user,
	}, nil
}
