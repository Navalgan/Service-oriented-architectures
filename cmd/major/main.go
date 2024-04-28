package main

import (
	"Service-oriented-architectures/internal/major"
	"github.com/gorilla/mux"

	"log"
	"net/http"
)

func main() {
	server := mux.NewRouter()

	service, err := major.NewService()
	if err != nil {
		log.Fatal(err)
	}

	server.HandleFunc("/user/join", service.UserJoin)
	server.HandleFunc("/user/auth", service.UserAuth)
	server.HandleFunc("/user/update", service.UserUpdate)
	server.HandleFunc("/post/create", service.CreatePost)
	server.HandleFunc("/post/{postId}", service.GetPostById)
	server.HandleFunc("/posts/{login}", service.GetPostsByLogin)
	server.HandleFunc("/post/{postId}/update", service.UpdatePost)
	server.HandleFunc("/post/{postId}/delete", service.DeletePost)

	log.Printf("Major service started")

	log.Fatal(http.ListenAndServe(":8080", server))
}
