package statistic

import (
	"Service-oriented-architectures/internal/statistic/grpc"

	"google.golang.org/grpc"
)

type Service struct {
	GRPCServer *grpc.Server
}

func NewService() (*Service, error) {
	gRPCServer := grpc.NewServer()

	if err := grpcstatistic.Register(gRPCServer); err != nil {
		return nil, err
	}

	return &Service{
		GRPCServer: gRPCServer,
	}, nil
}
