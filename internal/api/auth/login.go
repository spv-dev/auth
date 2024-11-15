package auth

import (
	"context"
	"errors"

	"github.com/spv-dev/auth/internal/model"
	"github.com/spv-dev/auth/internal/utils"
	authDesc "github.com/spv-dev/auth/pkg/auth_v1"
)

func (s *Server) Login(ctx context.Context, req *authDesc.LoginRequest) (*authDesc.LoginResponse, error) {
	refreshToken, err := utils.GenerateToken(model.TokenUserInfo{
		Username: req.GetUsername(),
		Role:     "ADMIN",
	},
		[]byte(refreshTokenSecretKey),
		refreshTokenExpiration,
	)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &authDesc.LoginResponse{
		RefreshToken: refreshToken,
	}, nil
}
