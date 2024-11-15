package auth

import (
	authDesc "github.com/spv-dev/auth/pkg/auth_v1"
)

// Server структура сервера
type Server struct {
	authDesc.UnimplementedAuthV1Server
}

// NewServer конструктор сервера
func NewServer() *Server {
	return &Server{}
}
