package access

import (
	"context"
	"errors"
	"strings"

	"github.com/spv-dev/auth/internal/model"
	"github.com/spv-dev/auth/internal/utils"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"

	descAccess "github.com/spv-dev/auth/pkg/access_v1"
)

var accessibleRoles map[string]string

const (
	authPrefix = "Bearer "
)

// Check метод проверки доступа к ресурсу
func (s *Server) Check(ctx context.Context, req *descAccess.CheckRequest) (*emptypb.Empty, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("metadata is not provided")
	}

	authHeader, ok := md["authorization"]
	if !ok || len(authHeader) == 0 {
		return nil, errors.New("authorization header is not provided")
	}

	if !strings.HasPrefix(authHeader[0], authPrefix) {
		return nil, errors.New("invalid authorization prefix format")
	}

	accessToken := strings.TrimPrefix(authHeader[0], authPrefix)

	claims, err := utils.VerifyToken(accessToken, []byte(s.config.GetAccessSecret()))
	if err != nil {
		return nil, errors.New("access token is invalid")
	}

	accessibleMap, err := s.accessibleRoles(ctx)
	if err != nil {
		return nil, errors.New("failed to get accessible roles")
	}

	role, ok := accessibleMap[req.GetEndpoint()]
	if !ok {
		return &emptypb.Empty{}, nil
	}

	if role == claims.Role {
		return &emptypb.Empty{}, nil
	}

	return nil, errors.New("access denied")
}

func (s *Server) accessibleRoles(_ context.Context) (map[string]string, error) {
	if accessibleRoles == nil {
		accessibleRoles = make(map[string]string)

		accessibleRoles[model.ExamplePath] = "ADMIN"
	}

	return accessibleRoles, nil
}
