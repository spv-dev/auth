package user

import (
	"context"

	"github.com/spv-dev/auth/internal/model"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *serv) UpdateUser(ctx context.Context, id int64, info *model.UpdateUserInfo) (*emptypb.Empty, error) {
	_, err := s.userRepository.UpdateUser(ctx, id, info)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
