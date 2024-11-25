package access

import (
	"context"
	"errors"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/spv-dev/auth/internal/constants"
	"github.com/spv-dev/auth/internal/model"
	serviceerror "github.com/spv-dev/auth/internal/service_error"
	"github.com/spv-dev/auth/internal/utils"
	descAccess "github.com/spv-dev/auth/pkg/access_v1"
)

var accessibleRoles map[string]string

// Check метод проверки доступа к ресурсу
func (s *Server) Check(ctx context.Context, req *descAccess.CheckRequest) (*emptypb.Empty, error) {
	accessToken, err := utils.GetAccessToken(ctx)
	if err != nil {
		return nil, err
	}

	claims, err := utils.VerifyToken(accessToken, []byte(s.config.GetAccessSecret()))
	if err != nil {
		return nil, errors.New(serviceerror.AccessTokenIsInvalid)
	}

	accessibleMap, err := s.accessibleRoles(ctx)
	if err != nil {
		return nil, errors.New(serviceerror.FailedToGetAccessibleRoles)
	}

	role, ok := accessibleMap[req.GetEndpoint()]
	if !ok {
		return &emptypb.Empty{}, nil
	}

	if role == claims.Role {
		return &emptypb.Empty{}, nil
	}

	return nil, errors.New(serviceerror.AccessDenied)
}

func (s *Server) accessibleRoles(_ context.Context) (map[string]string, error) {
	if accessibleRoles == nil {
		accessibleRoles = make(map[string]string)

		accessibleRoles[model.ExamplePath] = constants.ExampleRole
	}

	return accessibleRoles, nil
}
