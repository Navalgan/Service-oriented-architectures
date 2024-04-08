package grpctask

import (
	"Service-oriented-architectures/internal/common/gen/go"
	"Service-oriented-architectures/internal/task/storage"

	"context"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
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

func (s *ServiceAPI) GetPostById(ctx context.Context, req *task_v1.PostIdRequest) (*task_v1.PostResponse, error) {
	post, err := s.DB.GetPostById(req)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (s *ServiceAPI) GetPostsByLogin(ctx context.Context, req *task_v1.LoginRequest) (*task_v1.PostsResponse, error) {
	posts, err := s.DB.GetPostsByLogin(req)
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
