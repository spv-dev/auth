package main

import (
	"context"
	"log"

	"github.com/spv-dev/auth/internal/app"
)

func main() {
	ctx := context.Background()

	a, err := app.NewApp(ctx)
	if err != nil {
		log.Fatalf("failed to init app: %v", err.Error())
	}

	err = a.Run()
	if err != nil {
		log.Fatalf("Failed to run app: %v", err.Error())
	}
}
