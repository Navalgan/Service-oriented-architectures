package main

import (
	"Service-oriented-architectures/internal/posts"

	"log"
	"net"
)

func main() {
	service, err := posts.NewService()
	if err != nil {
		log.Fatal(err)
	}

	l, err := net.Listen("tcp", ":9090")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Task service started")

	log.Fatal(service.GRPCServer.Serve(l))
}
