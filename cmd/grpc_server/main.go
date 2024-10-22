package main

import (
	"context"
	"log"
	"net"

	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	userAPI "github.com/spv-dev/auth/internal/api/user"
	config "github.com/spv-dev/auth/internal/config"
	env "github.com/spv-dev/auth/internal/config/env"
	uRepository "github.com/spv-dev/auth/internal/repository/user"
	uService "github.com/spv-dev/auth/internal/service/user"
	desc "github.com/spv-dev/auth/pkg/user_v1"
)

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
	desc.RegisterAuthV1Server(s, userAPI.NewServer(userSrv))

	log.Printf("server listening at %s", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
