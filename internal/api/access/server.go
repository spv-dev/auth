package access

import (
	accessDesc "github.com/spv-dev/auth/pkg/access_v1"
)

// Server структура сервера
type Server struct {
	accessDesc.UnimplementedAccessV1Server
}

// NewServer конструктор сервера
func NewServer() *Server {
	return &Server{}
}
