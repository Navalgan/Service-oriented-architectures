package task

import (
	"Service-oriented-architectures/internal/task/grpc"
	"Service-oriented-architectures/internal/task/storage"

	"google.golang.org/grpc"
)

type Service struct {
	DB         *storage.DataBase
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
