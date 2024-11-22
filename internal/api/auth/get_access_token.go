package auth

import (
	"context"
	"errors"

	"github.com/spv-dev/auth/internal/model"
	"github.com/spv-dev/auth/internal/utils"
	authDesc "github.com/spv-dev/auth/pkg/auth_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetAccessToken метод для получения access token
func (s *Server) GetAccessToken(_ context.Context, req *authDesc.GetAccessTokenRequest) (*authDesc.GetAccessTokenResponse, error) {
	claims, err := utils.VerifyToken(req.GetRefreshToken(), []byte(s.config.GetRefreshSecret()))
	if err != nil {
		return nil, status.Errorf(codes.Aborted, "invalid refresh token")
	}

	accessToken, err := utils.GenerateToken(model.AuthUserInfo{
		Username: claims.Username,
		Role:     "ADMIN",
	},
		[]byte(s.config.GetAccessSecret()),
		s.config.GetAccessExpiration(),
	)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &authDesc.GetAccessTokenResponse{
		AccessToken: accessToken,
	}, nil
}
