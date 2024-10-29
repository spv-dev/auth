package user

import (
	"context"
	"log"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/spv-dev/auth/internal/converter"
	desc "github.com/spv-dev/auth/pkg/user_v1"
)

// UpdateUser изменяет информацию о пользователе
func (s *Server) UpdateUser(ctx context.Context, req *desc.UpdateUserRequest) (*emptypb.Empty, error) {

	err := s.userService.UpdateUser(ctx, req.GetId(), (converter.ToUpdateUserInfoFromDesc(req.GetInfo())))
	if err != nil {
		return nil, err
	}

	log.Printf("updated user user with id: %d", req.GetId())

	return nil, nil
}
