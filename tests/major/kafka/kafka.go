package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/IBM/sarama"
	"github.com/google/uuid"
)

func main() {
	producer, err := sarama.NewSyncProducer([]string{"localhost:29092"}, nil)
	if err != nil {
		log.Fatal(err)
	}

	postID := uuid.NewString()
	authorID := uuid.NewString()
	userID := uuid.NewString()
	date := time.Now().Unix()

	msg := &sarama.ProducerMessage{
		Topic: "likes",
		Key:   sarama.StringEncoder(uuid.New().String()),
		Value: sarama.StringEncoder(postID + "," + authorID + "," + userID + "," + strconv.FormatInt(date, 10)),
	}

	_, _, err = producer.SendMessage(msg)
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(time.Second * 5)

	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{"localhost:9000"},
	})

	row := conn.QueryRow(context.Background(), "SELECT count(distinct post_id, user_id) FROM Likes WHERE post_id=?;", postID)

	var count uint64
	if err := row.Scan(&count); err != nil {
		log.Fatal(err)
	}

	if count != 1 {
		log.Fatal("FAIL")
	}

	fmt.Println("OK")
}
