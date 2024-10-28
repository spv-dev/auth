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
	id := req.GetId()
	userInfo := converter.ToUpdateUserInfoFromDesc(req.GetInfo())
	err := s.userService.UpdateUser(ctx, id, &userInfo)
	if err != nil {
		return nil, err
	}

	log.Printf("updated user user with id: %d", id)

	return nil, nil
}
