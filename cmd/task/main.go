package main

import (
	"Service-oriented-architectures/internal/task"

	"log"
	"net"
)

func main() {
	service := task.NewService()

	l, err := net.Listen("tcp", ":9090")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Task service started")

	log.Fatal(service.GRPCServer.Serve(l))
}
