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

// GetRefreshToken метод для получения refresh token
func (s *Server) GetRefreshToken(_ context.Context, req *authDesc.GetRefreshTokenRequest) (*authDesc.GetRefreshTokenResponse, error) {
	claims, err := utils.VerifyToken(req.GetRefreshToken(), []byte(s.config.GetRefreshSecret()))
	if err != nil {
		return nil, status.Errorf(codes.Aborted, "invalid refresh token")
	}

	refreshToken, err := utils.GenerateToken(model.AuthUserInfo{
		Username: claims.Username,
		Role:     "ADMIN",
	},
		[]byte(s.config.GetRefreshSecret()),
		s.config.GetRefreshExpiration(),
	)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &authDesc.GetRefreshTokenResponse{
		RefreshToken: refreshToken,
	}, nil
}
