package main

import (
	"context"
	"log"
	"net"

	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"

	config "github.com/spv-dev/auth/internal/config"
	env "github.com/spv-dev/auth/internal/config/env"
	"github.com/spv-dev/auth/internal/converter"
	uRepository "github.com/spv-dev/auth/internal/repository/user"
	"github.com/spv-dev/auth/internal/service"
	uService "github.com/spv-dev/auth/internal/service/user"
	desc "github.com/spv-dev/auth/pkg/user_v1"
)

type server struct {
	desc.UnimplementedAuthV1Server
	userService service.UserService
}

func main() {
	ctx := context.Background()
	err := config.Load(".env")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	grpcConfig, err := env.NewGRPCConfig()
	if err != nil {
		log.Fatalf("failed to get grpc config: %v", err)
	}

	pgConfig, err := env.NewPGConfig()
	if err != nil {
		log.Fatalf("failed to get pg config: %v", err)
	}

	lis, err := net.Listen("tcp", grpcConfig.Address())
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	pool, err := pgxpool.Connect(ctx, pgConfig.DSN())
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	userRepo := uRepository.NewRepository(pool)
	userSrv := uService.NewService(userRepo)

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterAuthV1Server(s, &server{userService: userSrv})

	log.Printf("server listening at %s", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// CreateUser создаёт нового пользователя
func (s *server) CreateUser(ctx context.Context, req *desc.CreateUserRequest) (*desc.CreateUserResponse, error) {
	id, err := s.userService.CreateUser(ctx, converter.ToUserInfoFromDesc(req.GetInfo()), req.Password)
	if err != nil {
		return nil, err
	}

	log.Printf("inserted user with id: %d", id)

	return &desc.CreateUserResponse{
		Id: id,
	}, nil
}

// GetUser получает информацию о пользователе по идентификатору
func (s *server) GetUser(ctx context.Context, req *desc.GetUserRequest) (*desc.GetUserResponse, error) {
	id := req.GetId()
	userObj, err := s.userService.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}

	log.Printf("get user by id: %d", id)

	return &desc.GetUserResponse{
		User: converter.ToUserFromService(userObj),
	}, nil
}

// UpdateUser изменяет информацию о пользователе
func (s *server) UpdateUser(ctx context.Context, req *desc.UpdateUserRequest) (*emptypb.Empty, error) {
	id := req.GetId()
	_, err := s.userService.UpdateUser(ctx, id, converter.ToUpdateUserInfoFromDesc(req.GetInfo()))
	if err != nil {
		return nil, err
	}

	log.Printf("updated user user with id: %d", id)

	return nil, nil
}

// DeleteUser удаляет пользователя по идентификатору
func (s *server) DeleteUser(ctx context.Context, req *desc.DeleteUserRequest) (*emptypb.Empty, error) {
	id := req.GetId()
	_, err := s.userService.DeleteUser(ctx, id)
	if err != nil {
		return nil, err
	}

	log.Printf("delete user with id: %d", id)

	return nil, nil
}
