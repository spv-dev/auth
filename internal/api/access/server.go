package access

import (
	"github.com/spv-dev/auth/internal/config"
	accessDesc "github.com/spv-dev/auth/pkg/access_v1"
)

// Server структура сервера
type Server struct {
	accessDesc.UnimplementedAccessV1Server
	config config.AuthConfig
}

// NewServer конструктор сервера
func NewServer(config config.AuthConfig) *Server {
	return &Server{
		config: config,
	}
}
