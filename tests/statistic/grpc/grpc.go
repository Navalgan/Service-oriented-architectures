package main

import (
	"Service-oriented-architectures/internal/common/gen/go/statistic/proto"
	"fmt"
	"github.com/google/uuid"
	"time"

	"context"
	"log"

	"github.com/ClickHouse/clickhouse-go/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{"localhost:9000"},
	})

	postID := uuid.NewString()
	authorID := uuid.NewString()
	userID := uuid.NewString()
	date := time.Now().Unix()

	err = conn.Exec(context.Background(), "INSERT INTO Likes (post_id, author_id, user_id, date) VALUES (?, ?, ?, ?)", postID, authorID, userID, date)
	if err != nil {
		log.Fatal(err)
	}

	err = conn.Exec(context.Background(), "INSERT INTO Views (post_id, author_id, user_id, date) VALUES (?, ?, ?, ?)", postID, authorID, userID, date)
	if err != nil {
		log.Fatal(err)
	}

	conStatistic, err := grpc.Dial("localhost:7070", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	gRPCStatisticClient := statistic_v1.NewStatisticClient(conStatistic)

	resp, err := gRPCStatisticClient.GetPostStatByID(context.Background(), &statistic_v1.PostIDRequest{PostID: postID})

	if resp.Likes != 1 && resp.Views != 1 {
		log.Fatal("FAIL")
	}

	fmt.Println("OK")
}
