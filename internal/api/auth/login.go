package auth

import (
	"context"
	"errors"

	"github.com/spv-dev/auth/internal/model"
	"github.com/spv-dev/auth/internal/utils"
	authDesc "github.com/spv-dev/auth/pkg/auth_v1"
)

// Login метод для входа в систему
func (s *Server) Login(_ context.Context, req *authDesc.LoginRequest) (*authDesc.LoginResponse, error) {
	refreshToken, err := utils.GenerateToken(model.AuthUserInfo{
		Username: req.GetUsername(),
		Role:     "ADMIN",
	},
		[]byte(s.config.GetRefreshSecret()),
		s.config.GetRefreshExpiration(),
	)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &authDesc.LoginResponse{
		RefreshToken: refreshToken,
	}, nil
}
