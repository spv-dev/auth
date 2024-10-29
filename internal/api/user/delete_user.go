package user

import (
	"context"
	"log"

	"google.golang.org/protobuf/types/known/emptypb"

	desc "github.com/spv-dev/auth/pkg/user_v1"
)

// DeleteUser удаляет пользователя по идентификатору
func (s *Server) DeleteUser(ctx context.Context, req *desc.DeleteUserRequest) (*emptypb.Empty, error) {
	err := s.userService.DeleteUser(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	log.Printf("delete user with id: %d", req.GetId())

	return nil, nil
}
