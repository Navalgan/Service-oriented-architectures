package main

import (
	"Service-oriented-architectures/internal/major"

	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	service, err := major.NewService()
	if err != nil {
		log.Fatal(err)
	}

	mux.HandleFunc("/user/join", service.UserJoin)
	mux.HandleFunc("/user/auth", service.UserAuth)
	mux.HandleFunc("/user/update", service.UserUpdate)

	log.Printf("Major service started")

	log.Fatal(http.ListenAndServe(":8080", mux))
}
