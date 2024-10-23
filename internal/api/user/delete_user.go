package user

import (
	"context"
	"log"

	"google.golang.org/protobuf/types/known/emptypb"

	desc "github.com/spv-dev/auth/pkg/user_v1"
)

// DeleteUser удаляет пользователя по идентификатору
func (s *Server) DeleteUser(ctx context.Context, req *desc.DeleteUserRequest) (*emptypb.Empty, error) {
	id := req.GetId()
	_, err := s.userService.DeleteUser(ctx, id)
	if err != nil {
		return nil, err
	}

	log.Printf("delete user with id: %d", id)

	return nil, nil
}
