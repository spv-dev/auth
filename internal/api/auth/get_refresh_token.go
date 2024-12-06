package auth

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/spv-dev/auth/internal/constants"
	"github.com/spv-dev/auth/internal/model"
	serviceerror "github.com/spv-dev/auth/internal/service_error"
	"github.com/spv-dev/auth/internal/utils"
	authDesc "github.com/spv-dev/auth/pkg/auth_v1"
)

// GetRefreshToken метод для получения refresh token
func (s *Server) GetRefreshToken(_ context.Context, req *authDesc.GetRefreshTokenRequest) (*authDesc.GetRefreshTokenResponse, error) {
	claims, err := utils.VerifyToken(req.GetRefreshToken(), []byte(s.config.GetRefreshSecret()))
	if err != nil {
		return nil, status.Errorf(codes.Aborted, serviceerror.InvalidRefreshToken)
	}

	refreshToken, err := utils.GenerateToken(model.AuthUserInfo{
		Username: claims.Username,
		Role:     constants.ExampleRole,
	},
		[]byte(s.config.GetRefreshSecret()),
		s.config.GetRefreshExpiration(),
	)
	if err != nil {
		return nil, errors.New(serviceerror.FailedToGenerateToken)
	}

	return &authDesc.GetRefreshTokenResponse{
		RefreshToken: refreshToken,
	}, nil
}
