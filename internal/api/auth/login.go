package auth

import (
	"context"
	"errors"

	"github.com/spv-dev/auth/internal/constants"
	"github.com/spv-dev/auth/internal/model"
	serviceerror "github.com/spv-dev/auth/internal/service_error"
	"github.com/spv-dev/auth/internal/utils"
	authDesc "github.com/spv-dev/auth/pkg/auth_v1"
)

// Login метод для входа в систему
func (s *Server) Login(_ context.Context, req *authDesc.LoginRequest) (*authDesc.LoginResponse, error) {
	refreshToken, err := utils.GenerateToken(model.AuthUserInfo{
		Username: req.GetUsername(),
		Role:     constants.ExampleRole,
	},
		[]byte(s.config.GetRefreshSecret()),
		s.config.GetRefreshExpiration(),
	)
	if err != nil {
		return nil, errors.New(serviceerror.FailedToGenerateRefreshToken)
	}

	accessToken, err := utils.GenerateToken(model.AuthUserInfo{
		Username: req.GetUsername(),
		Role:     constants.ExampleRole,
	},
		[]byte(s.config.GetAccessSecret()),
		s.config.GetAccessExpiration(),
	)
	if err != nil {
		return nil, errors.New(serviceerror.FailedToGenerateAccessToken)
	}

	return &authDesc.LoginResponse{
		RefreshToken: refreshToken,
		AccessToken:  accessToken,
	}, nil
}
