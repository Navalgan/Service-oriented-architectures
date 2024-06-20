package main

import (
	"Service-oriented-architectures/internal/statistic"

	"log"
	"net"
)

func main() {
	service, err := statistic.NewService()
	if err != nil {
		log.Fatal(err)
	}

	l, err := net.Listen("tcp", ":7070")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Statistic service started")

	log.Fatal(service.GRPCServer.Serve(l))
}
