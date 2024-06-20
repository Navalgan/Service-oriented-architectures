package main

import (
	"Service-oriented-architectures/internal/major"

	"context"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/spf13/pflag"
)

var (
	flagJwtKey = pflag.String("jwt-key", "", "the jwt key")
)

func main() {
	pflag.Parse()

	ctx := context.Background()

	service, err := major.NewService(*flagJwtKey, ctx)
	if err != nil {
		log.Fatal(err)
	}

	server := mux.NewRouter()

	server.HandleFunc("/user/join", service.UserJoin)
	server.HandleFunc("/user/auth", service.UserAuth)

	server.HandleFunc("/user/update", service.UserUpdate)
	server.HandleFunc("/post/create", service.CreatePost)
	server.HandleFunc("/post/{postId}", service.GetPostByID)
	server.HandleFunc("/post/{postId}/stat", service.GetPostStatByID)
	server.HandleFunc("/posts/{login}", service.GetPostsByUser)
	server.HandleFunc("/post/{postId}/like", service.LikePost)
	server.HandleFunc("/post/{postId}/view", service.ViewPost)
	server.HandleFunc("/post/{postId}/update", service.UpdatePost)
	server.HandleFunc("/post/{postId}/delete", service.DeletePost)

	server.HandleFunc("/top/users", service.GetTopUsers)
	server.HandleFunc("/top/posts", service.GetTopPosts)

	log.Printf("Major service started")

	log.Fatal(http.ListenAndServe(":8080", server))
}
