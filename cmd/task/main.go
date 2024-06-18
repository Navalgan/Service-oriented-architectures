package main

import (
	"Service-oriented-architectures/internal/task"

	"log"
	"net"
)

func main() {
	service, err := task.NewService()
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
