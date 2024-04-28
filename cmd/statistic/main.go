package main

import (
	"Service-oriented-architectures/internal/statistic"

	"context"
	"log"
)

func main() {
	ctx := context.Background()

	service, err := statistic.NewService(ctx)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Statistic service started")

	service.Run(ctx)
}
