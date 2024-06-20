package grpcstatistic

import (
	"Service-oriented-architectures/internal/common/gen/go/statistic/proto"
	"Service-oriented-architectures/internal/statistic/storage"

	"context"
	"log"

	"google.golang.org/grpc"
)

type ServiceAPI struct {
	statistic_v1.UnimplementedStatisticServer
	DB *storage.DataBase
}

func Register(gRPC *grpc.Server) error {
	ctx := context.Background()

	db, err := storage.NewDataBase(ctx)
	if err != nil {
		log.Print(err)
		return err
	}

	statistic_v1.RegisterStatisticServer(gRPC, &ServiceAPI{DB: db})
	return nil
}

func (s *ServiceAPI) GetPostStatByID(ctx context.Context, req *statistic_v1.PostIDRequest) (*statistic_v1.PostStatResponse, error) {
	likesCount, err := s.DB.GetLikesCount(ctx, req.PostID)
	if err != nil {
		return nil, err
	}

	viewsCount, err := s.DB.GetViewsCount(ctx, req.PostID)
	if err != nil {
		return nil, err
	}

	return &statistic_v1.PostStatResponse{Likes: likesCount, Views: viewsCount}, nil
}

func (s *ServiceAPI) GetTopPosts(ctx context.Context, req *statistic_v1.TopPostsRequest) (*statistic_v1.TopPostsResponse, error) {
	return s.DB.GetTopPosts(ctx, req.OrderBy)
}

func (s *ServiceAPI) GetTopUsers(ctx context.Context, _ *statistic_v1.TopUsersRequest) (*statistic_v1.TopUsersResponse, error) {
	return s.DB.GetTopUsers(ctx)
}
