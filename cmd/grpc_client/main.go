package main

import (
	"context"
	"log"
	"time"

	desc "github.com/spv-dev/auth/pkg/user_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const host = "localhost:50051"

func main() {
	conn, err := grpc.NewClient(host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("can't to create client: %v", err)
	}

	defer func() {
		if err := conn.Close(); err != nil {
			log.Fatalf("failed when close connection: %v", err)
		}
	}()

	c := desc.NewAuthV1Client(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	respCreate, err := c.Create(ctx, &desc.CreateRequest{
		Info: &desc.UserInfo{
			Name:  "Max",
			Email: "max@gmail.com",
			Role:  desc.Roles_ADMIN,
		},
		Password:        "pass1",
		PasswordConfirm: "pass1",
	})
	if err != nil {
		log.Fatalf("failed when create: %v", err)
	}
	log.Printf("Create user: \n%v", respCreate.GetId())

	respUpdate, err := c.Update(ctx, &desc.UpdateRequest{})
	if err != nil {
		log.Fatalf("failed when update: %v", err)
	}
	log.Printf("Update user ok: \n%v", respUpdate)

	respDelete, err := c.Delete(ctx, &desc.DeleteRequest{})
	if err != nil {
		log.Fatalf("failed when delete: %v", err)
	}
	log.Printf("Delete user ok: \n%v", respDelete)

	respGet, err := c.Get(ctx, &desc.GetRequest{})
	if err != nil {
		log.Fatalf("failed when get: %v", err)
	}
	log.Printf("Get user ok: \n%v", respGet)
}
