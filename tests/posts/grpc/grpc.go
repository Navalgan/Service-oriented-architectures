package main

import (
	"Service-oriented-architectures/internal/common/gen/go/posts/proto"

	"context"
	"fmt"
	"github.com/google/uuid"
	"strings"

	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conPosts, err := grpc.Dial("localhost:9090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	gRPCPostsClient := posts_v1.NewPostsClient(conPosts)

	userID := uuid.NewString()
	text := "Posts test"

	createdPost, err := gRPCPostsClient.CreatePost(context.Background(), &posts_v1.PostRequest{UserID: userID, Text: text})
	if err != nil {
		fmt.Println("FAIL")
		log.Fatal(err)
	}

	resp, err := gRPCPostsClient.GetPostByID(context.Background(), &posts_v1.PostIDRequest{PostID: createdPost.PostID})
	if err != nil {
		fmt.Println("FAIL")
		log.Fatal(err)
	}

	if strings.Compare(createdPost.PostID, resp.PostID) != 0 && strings.Compare(createdPost.AuthorID, resp.AuthorID) != 0 &&
		strings.Compare(createdPost.Text, resp.Text) != 0 && createdPost.Date != resp.Date {
		log.Fatal("FAIL")
	}

	fmt.Println("OK")
}
