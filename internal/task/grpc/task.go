package grpctask

import (
	"Service-oriented-architectures/internal/common/gen/go/task/proto"
	"Service-oriented-architectures/internal/task/storage"

	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ServiceAPI struct {
	task_v1.UnimplementedTaskServer
	DB *storage.DataBase
}

func Register(gRPC *grpc.Server) error {
	db, err := storage.NewDataBase()
	if err != nil {
		log.Print(err)
		return err
	}

	task_v1.RegisterTaskServer(gRPC, &ServiceAPI{DB: db})
	return nil
}

func (s *ServiceAPI) CreatePost(ctx context.Context, req *task_v1.PostRequest) (*task_v1.PostResponse, error) {
	postId, err := s.DB.CreatePost(req)
	if err != nil {
		return nil, err
	}
	return postId, nil
}

func (s *ServiceAPI) GetPostByID(ctx context.Context, req *task_v1.PostIDRequest) (*task_v1.PostResponse, error) {
	post, err := s.DB.GetPostByID(req)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (s *ServiceAPI) GetPostsByUser(ctx context.Context, req *task_v1.UserRequest) (*task_v1.PostsResponse, error) {
	posts, err := s.DB.GetPostsByUser(req)
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (s *ServiceAPI) UpdatePost(ctx context.Context, req *task_v1.UpdatePostRequest) (*emptypb.Empty, error) {
	if err := s.DB.UpdatePost(req); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *ServiceAPI) DeletePost(ctx context.Context, req *task_v1.DeletePostRequest) (*emptypb.Empty, error) {
	if err := s.DB.DeletePost(req); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
