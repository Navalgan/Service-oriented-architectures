package storage

import (
	"Service-oriented-architectures/internal/common/gen/go/statistic/proto"

	"context"
	"fmt"
	"log"

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

	err = conn.Exec(ctx, "CREATE TABLE IF NOT EXISTS UsersInteraction (author_id UUID, user_id UUID, date DateTime('Europe/Moscow')) Engine = MergeTree ORDER BY author_id;")
	if err != nil {
		return nil, err
	}
	log.Printf("Target table for likes created")

	err = conn.Exec(ctx, "CREATE TABLE IF NOT EXISTS Likes (post_id UUID, author_id UUID, user_id UUID, date DateTime('Europe/Moscow')) Engine = MergeTree ORDER BY post_id;")
	if err != nil {
		return nil, err
	}
	log.Printf("Target table for likes created")

	err = conn.Exec(ctx, "CREATE TABLE IF NOT EXISTS LikesQueue (post_id UUID, author_id UUID, user_id UUID, date DateTime('Europe/Moscow')) Engine = Kafka SETTINGS kafka_broker_list = 'kafka:9092', kafka_topic_list = 'likes', kafka_group_name = 'likes_consumer_group1', kafka_format = 'CSV', kafka_max_block_size = 1048576;")
	if err != nil {
		return nil, err
	}
	log.Printf("LikesQueue table created")

	err = conn.Exec(ctx, "CREATE MATERIALIZED VIEW IF NOT EXISTS UsersInteractionMV TO UsersInteraction AS SELECT author_id, user_id, date FROM LikesQueue;")
	if err != nil {
		return nil, err
	}
	log.Printf("Materialized view for likes created")

	err = conn.Exec(ctx, "CREATE MATERIALIZED VIEW IF NOT EXISTS LikesQueueMV TO Likes AS SELECT post_id, author_id, user_id, date FROM LikesQueue;")
	if err != nil {
		return nil, err
	}
	log.Printf("Materialized view for likes created")

	err = conn.Exec(ctx, "CREATE TABLE IF NOT EXISTS Views (post_id UUID, author_id UUID, user_id UUID, date DateTime('Europe/Moscow')) Engine = MergeTree ORDER BY post_id;")
	if err != nil {
		return nil, err
	}
	log.Printf("Target table for views created")

	err = conn.Exec(ctx, "CREATE TABLE IF NOT EXISTS ViewsQueue (post_id UUID, author_id UUID, user_id UUID, date DateTime('Europe/Moscow')) Engine = Kafka SETTINGS kafka_broker_list = 'kafka:9092', kafka_topic_list = 'views', kafka_group_name = 'views_consumer_group1 ', kafka_format = 'CSV', kafka_max_block_size = 1048576;")
	if err != nil {
		return nil, err
	}
	log.Printf("Views table created")

	err = conn.Exec(ctx, "CREATE MATERIALIZED VIEW IF NOT EXISTS ViewsQueueMv TO Views AS SELECT post_id, author_id, user_id, date FROM ViewsQueue;")
	if err != nil {
		return nil, err
	}
	log.Printf("Materialized view for views created")

	return &DataBase{Conn: conn}, nil
}

func (db *DataBase) GetLikesCount(ctx context.Context, postID string) (uint64, error) {
	const query = "SELECT count(distinct post_id, user_id) FROM Likes WHERE post_id=?;"

	row := db.Conn.QueryRow(ctx, query, postID)

	var count uint64
	if err := row.Scan(&count); err != nil {
		log.Println("like count:", err)
		return 0, err
	}

	return count, nil
}

func (db *DataBase) GetViewsCount(ctx context.Context, postID string) (uint64, error) {
	const query = "SELECT count(distinct post_id, user_id) FROM Views WHERE post_id=?;"

	row := db.Conn.QueryRow(ctx, query, postID)

	var count uint64
	if err := row.Scan(&count); err != nil {
		log.Println("view count:", err)
		return 0, err
	}

	return count, nil
}

func (db *DataBase) GetTopPosts(ctx context.Context, tableName string) (*statistic_v1.TopPostsResponse, error) {
	query := "SELECT t.post_id, t.author_id, count(distinct post_id, user_id) AS field_count FROM " + tableName + " t GROUP BY post_id, author_id ORDER BY field_count DESC LIMIT 5;"

	rows, err := db.Conn.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	resp := make([]*statistic_v1.PostStatistic, 0)
	for rows.Next() {
		iter := &statistic_v1.PostStatistic{}
		err = rows.Scan(&iter.PostID, &iter.AuthorID, &iter.Count)
		if err != nil {
			log.Println("row scan:", err)
			return nil, err
		}
		resp = append(resp, iter)
	}

	log.Printf("Top of %d posts", len(resp))

	return &statistic_v1.TopPostsResponse{Posts: resp}, nil
}

func (db *DataBase) GetTopUsers(ctx context.Context) (*statistic_v1.TopUsersResponse, error) {
	const query = "SELECT author_id, count(distinct author_id, user_id) AS likes_count FROM UsersInteraction GROUP BY author_id ORDER BY likes_count DESC LIMIT 3;"

	rows, err := db.Conn.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	resp := make([]*statistic_v1.UserStatistic, 0)
	for rows.Next() {
		iter := &statistic_v1.UserStatistic{}
		err = rows.Scan(&iter.UserID, &iter.Count)
		if err != nil {
			log.Println("row scan:", err)
			return nil, err
		}
		resp = append(resp, iter)
	}

	log.Printf("Top of %d uesrs", len(resp))

	return &statistic_v1.TopUsersResponse{Users: resp}, nil
}
