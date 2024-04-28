package storage

import (
	"Service-oriented-architectures/internal/common"
	"context"
	"fmt"
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

	err = conn.Exec(ctx, "CREATE TABLE IF NOT EXISTS statistics (post_id UUID, user_id String, operation UInt8) Engine = Memory")
	if err != nil {
		return nil, err
	}

	return &DataBase{Conn: conn}, nil
}

func (db *DataBase) SetLike(ctx context.Context, postID, userID string) error {
	const query = "INSERT INTO statistics (post_id, user_id, operation) VALUES (?, ?, ?)"

	return db.Conn.Exec(ctx, query, postID, userID, common.Like, postID, userID, common.Like)
}

func (db *DataBase) SetView(ctx context.Context, postID, userID string) error {
	const query = "INSERT INTO statistics (post_id, user_id, operation) VALUES (?, ?, ?)"

	return db.Conn.Exec(ctx, query, postID, userID, common.View)
}
