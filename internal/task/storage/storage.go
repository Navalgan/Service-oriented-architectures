package storage

import (
	"Service-oriented-architectures/internal/common/gen/go/task/proto"

	"log"
	"time"

	"github.com/gocql/gocql"
	"github.com/google/uuid"
)

type DataBase struct {
	s *gocql.Session
}

func NewDataBase() (*DataBase, error) {
	time.Sleep(10 * time.Second)

	cluster := gocql.NewCluster("cassandra1", "cassandra2", "cassandra3")
	cluster.Consistency = gocql.Quorum
	cluster.ConnectTimeout = time.Second * 10

	session, err := cluster.CreateSession()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if err = session.Query(
		"CREATE KEYSPACE IF NOT EXISTS SocialNetwork WITH REPLICATION = {'class': 'NetworkTopologyStrategy'}",
	).Exec(); err != nil {
		log.Println(err)
		return nil, err
	}

	log.Print("KEYSPACE ready")

	const query = "CREATE TABLE IF NOT EXISTS SocialNetwork.Posts (post_id uuid, author_id uuid, last_update bigint, post_text text, PRIMARY KEY (post_id, author_id));"

	if err = session.Query(query).Exec(); err != nil {
		log.Println(err)
		return nil, err
	}

	log.Print("Cassandra ready")

	return &DataBase{
		session,
	}, nil
}

func (db *DataBase) CreatePost(req *task_v1.PostRequest) (*task_v1.PostResponse, error) {
	const query = "INSERT INTO SocialNetwork.Posts (post_id, author_id, last_update, post_text) VALUES (?, ?, ?, ?)"

	postID := uuid.NewString()

	curTime := time.Now().UnixMilli()

	if err := db.s.Query(query, postID, req.UserID, curTime, req.Text).Exec(); err != nil {
		log.Println(err)
		return nil, err
	}

	return &task_v1.PostResponse{
		PostID:   postID,
		AuthorID: req.UserID,
		Date:     curTime,
		Text:     req.Text,
	}, nil
}

func (db *DataBase) GetPostByID(req *task_v1.PostIDRequest) (*task_v1.PostResponse, error) {
	const query = "SELECT post_id, author_id, last_update, post_text FROM SocialNetwork.Posts WHERE post_id=?"

	var post task_v1.PostResponse

	if err := db.s.Query(query, req.PostID).Scan(&post.PostID, &post.AuthorID, &post.Date, &post.Text); err != nil {
		log.Println(err)
		return nil, err
	}

	return &post, nil
}

func (db *DataBase) GetPostsByUser(req *task_v1.UserRequest) (*task_v1.PostsResponse, error) {
	const query = "SELECT post_id, author_id, last_update, post_text FROM SocialNetwork.Posts WHERE author_id=? ALLOW FILTERING"

	var posts task_v1.PostsResponse

	scanner := db.s.Query(query, req.UserID).Iter().Scanner()
	for scanner.Next() {
		var post task_v1.PostResponse
		err := scanner.Scan(&post.PostID, &post.AuthorID, &post.Date, &post.Text)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		posts.Posts = append(posts.Posts, &post)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return &posts, nil
}

func (db *DataBase) UpdatePost(req *task_v1.UpdatePostRequest) error {
	const query = "UPDATE SocialNetwork.Posts SET last_update=?, post_text=? WHERE post_id=? and author_id=?"

	if err := db.s.Query(query, time.Now().Unix(), req.Text, req.PostID, req.UserID).Exec(); err != nil {
		log.Println(err)
	}

	return nil
}

func (db *DataBase) DeletePost(req *task_v1.DeletePostRequest) error {
	const query = "DELETE FROM SocialNetwork.Posts WHERE post_id=? and author_id=?"

	if err := db.s.Query(query, req.PostID, req.UserID).Exec(); err != nil {
		log.Println(err)
	}

	return nil
}
