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

// GetAccessToken метод для получения access token
func (s *Server) GetAccessToken(_ context.Context, req *authDesc.GetAccessTokenRequest) (*authDesc.GetAccessTokenResponse, error) {
	claims, err := utils.VerifyToken(req.GetRefreshToken(), []byte(s.config.GetRefreshSecret()))
	if err != nil {
		return nil, status.Errorf(codes.Aborted, serviceerror.InvalidRefreshToken)
	}

	accessToken, err := utils.GenerateToken(model.AuthUserInfo{
		Username: claims.Username,
		Role:     constants.ExampleRole,
	},
		[]byte(s.config.GetAccessSecret()),
		s.config.GetAccessExpiration(),
	)
	if err != nil {
		return nil, errors.New(serviceerror.FailedToGenerateAccessToken)
	}

	return &authDesc.GetAccessTokenResponse{
		AccessToken: accessToken,
	}, nil
}
