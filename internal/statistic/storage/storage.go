package storage

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
)

type DataBase struct {
	Conn clickhouse.Conn
}

func NewDataBase(ctx context.Context) (*DataBase, error) {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{"clickhouse:9000"},
	})

	if err != nil {
		return nil, err
	}

	if err = conn.Ping(ctx); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			fmt.Printf("Exception [%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		}
		return nil, err
	}

	err = conn.Exec(ctx, "CREATE TABLE IF NOT EXISTS Likes (post_id UUID, user_id UUID, date DateTime('Europe/Moscow')) Engine = MergeTree ORDER BY tuple();")
	if err != nil {
		return nil, err
	}
	log.Printf("Target table for likes created")

	err = conn.Exec(ctx, "CREATE TABLE IF NOT EXISTS LikesQueue (post_id UUID, user_id UUID, date DateTime('Europe/Moscow')) Engine = Kafka SETTINGS kafka_broker_list = 'kafka:9092', kafka_topic_list = 'likes', kafka_group_name = 'likes_consumer_group1', kafka_format = 'CSV', kafka_max_block_size = 1048576;")
	if err != nil {
		return nil, err
	}
	log.Printf("LikesQueue table created")

	err = conn.Exec(ctx, "CREATE MATERIALIZED VIEW IF NOT EXISTS LikesQueueMv TO Likes AS SELECT post_id, user_id, date FROM LikesQueue;")
	if err != nil {
		return nil, err
	}
	log.Printf("Materialized view for likes created")

	err = conn.Exec(ctx, "CREATE TABLE IF NOT EXISTS Views (post_id UUID, user_id UUID, date DateTime('Europe/Moscow')) Engine = MergeTree ORDER BY tuple();")
	if err != nil {
		return nil, err
	}
	log.Printf("Target table for views created")

	err = conn.Exec(ctx, "CREATE TABLE IF NOT EXISTS ViewsQueue (post_id UUID, user_id UUID, date DateTime('Europe/Moscow')) Engine = Kafka SETTINGS kafka_broker_list = 'kafka:9092', kafka_topic_list = 'views', kafka_group_name = 'views_consumer_group1 ', kafka_format = 'CSV', kafka_max_block_size = 1048576;")
	if err != nil {
		return nil, err
	}
	log.Printf("Views table created")

	err = conn.Exec(ctx, "CREATE MATERIALIZED VIEW IF NOT EXISTS ViewsQueueMv TO Views AS SELECT post_id, user_id, date FROM ViewsQueue;")
	if err != nil {
		return nil, err
	}
	log.Printf("Materialized view for views created")

	return &DataBase{Conn: conn}, nil
}

func (db *DataBase) SetLike(ctx context.Context, postID, userID string) error {
	const query = "INSERT INTO Likes (post_id, user_id, time) VALUES (?, ?, ?)"

	return db.Conn.Exec(ctx, query, postID, userID, time.Now().UnixMilli())
}

func (db *DataBase) SetView(ctx context.Context, postID, userID string) error {
	const query = "INSERT INTO Views (post_id, user_id, time) VALUES (?, ?, ?)"

	return db.Conn.Exec(ctx, query, postID, userID, time.Now().UnixMilli())
}
