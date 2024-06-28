package posts

import (
	"Service-oriented-architectures/internal/posts/grpc"

	"google.golang.org/grpc"
)

type Service struct {
	GRPCServer *grpc.Server
}

func NewService() (*Service, error) {
	gRPCServer := grpc.NewServer()

	if err := grpctask.Register(gRPCServer); err != nil {
		return nil, err
	}

	return &Service{
		GRPCServer: gRPCServer,
	}, nil
}
