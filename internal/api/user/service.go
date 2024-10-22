package user

import (
	"github.com/spv-dev/auth/internal/service"
	desc "github.com/spv-dev/auth/pkg/user_v1"
)

type Server struct {
	desc.UnimplementedAuthV1Server
	userService service.UserService
}

func NewServer(userService service.UserService) *Server {
	return &Server{
		userService: userService,
	}
}
