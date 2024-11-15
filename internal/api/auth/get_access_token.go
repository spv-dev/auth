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

func (s *Server) GetAccessToken(ctx context.Context, req *authDesc.GetAccessTokenRequest) (*authDesc.GetAccessTokenResponse, error) {
	claims, err := utils.VerifyToken(req.GetRefreshToken(), []byte(refreshTokenSecretKey))
	if err != nil {
		return nil, status.Errorf(codes.Aborted, "invalid refresh token")
	}

	accessToken, err := utils.GenerateToken(model.TokenUserInfo{
		Username: claims.Username,
		Role:     "ADMIN",
	},
		[]byte(accessTokenSecretKey),
		accessTokenExpiration,
	)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &authDesc.GetAccessTokenResponse{
		AccessToken: accessToken,
	}, nil
}
