package utils

import (
	"context"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"google.golang.org/grpc/metadata"

	"github.com/spv-dev/auth/internal/model"
	serviceerror "github.com/spv-dev/auth/internal/service_error"
)

const (
	authPrefix = "Bearer "
)

// GenerateToken метод генерирует новый токен
func GenerateToken(info model.AuthUserInfo, secretKey []byte, duration time.Duration) (string, error) {
	claims := model.UserClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(duration).Unix(),
		},
		Username: info.Username,
		Role:     info.Role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(secretKey)
}

// VerifyToken метод проверяет пришедший токен
func VerifyToken(tokenStr string, secretKey []byte) (*model.UserClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&model.UserClaims{},
		func(t *jwt.Token) (interface{}, error) {
			_, ok := t.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, errors.Errorf(serviceerror.UnexpectedTokenSigningMethod)
			}

			return secretKey, nil
		},
	)

	if err != nil {
		return nil, errors.Errorf(serviceerror.InvalidToken, err.Error())
	}

	claims, ok := token.Claims.(*model.UserClaims)
	if !ok {
		return nil, errors.Errorf(serviceerror.InvalidTokenClaim)
	}

	return claims, nil
}

// GetAccessToken получение токена из хэдера
func GetAccessToken(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", errors.New(serviceerror.MetadataIsNotProvided)
	}

	authHeader, ok := md["authorization"]
	if !ok || len(authHeader) == 0 {
		return "", errors.New(serviceerror.AuthorizationHeaderIsNotProvided)
	}

	if !strings.HasPrefix(authHeader[0], authPrefix) {
		return "", errors.New(serviceerror.InvalidAuthorizationPrefixFormat)
	}

	accessToken := strings.TrimPrefix(authHeader[0], authPrefix)

	return accessToken, nil
}
