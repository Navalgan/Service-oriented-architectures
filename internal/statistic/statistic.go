package statistic

import (
	"Service-oriented-architectures/internal/task/grpc"
	"Service-oriented-architectures/internal/task/storage"

	"google.golang.org/grpc"
)

type Service struct {
	GRPCServer *grpc.Server
	DB         *storage.DataBase
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
