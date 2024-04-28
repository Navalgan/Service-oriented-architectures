package storage

import (
	"Service-oriented-architectures/internal/common/gen/go"

	"context"
	"log"
	"time"

	"github.com/gocql/gocql"
	"github.com/google/uuid"
)

type DataBase struct {
	s *gocql.Session
}

const (
	Base = 10
)

func NewDataBase() (*DataBase, error) {
	time.Sleep(10 * time.Second)

	cluster := gocql.NewCluster("cassandra1", "cassandra2", "cassandra3")
	cluster.Consistency = gocql.Quorum
	cluster.ProtoVersion = 4
	cluster.ConnectTimeout = time.Second * 10
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: "Username",
		Password: "Password",
	}

	session, err := cluster.CreateSession()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if err = session.Query(
		"CREATE KEYSPACE IF NOT EXISTS SocialNetwork WITH REPLICATION = {'class': 'NetworkTopologyStrategy'}",
	).Exec(); err != nil {
		log.Fatal(err)
		return nil, err
	}

	log.Print("KEYSPACE ready")

	const query = "CREATE TABLE IF NOT EXISTS SocialNetwork.Posts (postId uuid, lastUpdate text, author text, postText text, PRIMARY KEY (postId, author));"

	if err = session.Query(query).Exec(); err != nil {
		log.Println(err)
		return nil, err
	}

	log.Print("cassandra ready")

	return &DataBase{
		session,
	}, nil
}

func (db *DataBase) CreatePost(req *task_v1.PostRequest) (*task_v1.PostResponse, error) {
	const query = "INSERT INTO SocialNetwork.Posts (postId, lastUpdate, author, postText) VALUES (?, ?, ?, ?)"

	postId := uuid.NewString()

	curTime := time.Now().Format(time.RFC3339)

	if err := db.s.Query(query, postId, curTime, req.Login, req.Text).Exec(); err != nil {
		log.Println(err)
		return nil, err
	}

	return &task_v1.PostResponse{
		PostId: postId,
		Date:   curTime,
		Author: req.Login,
		Text:   req.Text,
	}, nil
}

func (db *DataBase) GetPostById(req *task_v1.PostIdRequest) (*task_v1.PostResponse, error) {
	const query = "SELECT * FROM SocialNetwork.Posts WHERE postId=?"

	var post task_v1.PostResponse
	if err := db.s.Query(query, req.PostId).WithContext(context.Background()).Scan(&post.PostId, &post.Date, &post.Author, &post.Text); err != nil {
		log.Println(err)
		return nil, err
	}

	return &post, nil
}

func (db *DataBase) GetPostsByLogin(req *task_v1.LoginRequest) (*task_v1.PostsResponse, error) {
	const query = "SELECT * FROM SocialNetwork.Posts WHERE author=? ALLOW FILTERING"

	var posts task_v1.PostsResponse

	scanner := db.s.Query(query, req.Login).Iter().Scanner()
	for scanner.Next() {
		var post task_v1.PostResponse
		err := scanner.Scan(&post.PostId, &post.Date, &post.Author, &post.Text)
		if err != nil {
			log.Fatal(err)
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
	const query = "UPDATE SocialNetwork.Posts SET lastUpdate=?, postText=? WHERE postId=? and author=?"

	if err := db.s.Query(query, time.Now().Format(time.RFC3339), req.Text, req.PostId, req.Login).Exec(); err != nil {
		log.Println(err)
	}

	return nil
}

func (db *DataBase) DeletePost(req *task_v1.DeletePostRequest) error {
	const query = "DELETE FROM SocialNetwork.Posts WHERE postId=? and author=?"

	if err := db.s.Query(query, req.PostId, req.Login).Exec(); err != nil {
		log.Println(err)
	}

	return nil
}
