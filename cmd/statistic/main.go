package main

import (
	"Service-oriented-architectures/internal/statistic"

	"log"
	"net"
)

func main() {
	service, err := statistic.NewService()
	if err != nil {
		log.Fatal("")
	}

	l, err := net.Listen("tcp", ":9090")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Statistic service started")

	log.Fatal(service.Serve(l))
}
