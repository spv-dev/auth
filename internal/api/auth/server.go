package auth

import (
	"github.com/spv-dev/auth/internal/config"
	authDesc "github.com/spv-dev/auth/pkg/auth_v1"
)

// Server структура сервера
type Server struct {
	authDesc.UnimplementedAuthV1Server
	config config.AuthConfig
}

// NewServer конструктор сервера
func NewServer(config config.AuthConfig) *Server {
	return &Server{
		config: config,
	}
}
