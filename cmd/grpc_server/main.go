package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"

	desc "github.com/spv-dev/auth/pkg/user_v1"
)

const grpcPort = 50051

type server struct {
	desc.UnimplementedAuthV1Server
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterAuthV1Server(s, &server{})

	log.Printf("server listening at %s", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// CreateUser created a new user
// Return id of created user
func (s *server) CreateUser(ctx context.Context, req *desc.CreateUserRequest) (*desc.CreateUserResponse, error) {
	log.Printf("Method: %s\nRequest: %v\nContext: %v\n", "CreateUser", req, ctx)
	return &desc.CreateUserResponse{}, nil
}

// GetUser gets user info by id
// Return info about user
func (s *server) GetUser(ctx context.Context, req *desc.GetUserRequest) (*desc.GetUserResponse, error) {
	log.Printf("Method: %s\nRequest: %v\nContext: %v\n", "GetUser", req, ctx)
	return &desc.GetUserResponse{}, nil
}

// UpdateUser changes info about user
func (s *server) UpdateUser(ctx context.Context, req *desc.UpdateUserRequest) (*emptypb.Empty, error) {
	log.Printf("Method: %s\nRequest: %v\nContext: %v\n", "UpdateUser", req, ctx)
	return nil, nil
}

// DeleteUser deletes user by id
func (s *server) DeleteUser(ctx context.Context, req *desc.DeleteUserRequest) (*emptypb.Empty, error) {
	log.Printf("Method: %s\nRequest: %v\nContext: %v\n", "DeleteUser", req, ctx)
	return nil, nil
}
