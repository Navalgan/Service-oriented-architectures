package statistic

import (
	"Service-oriented-architectures/internal/statistic/storage"

	"context"
	"log"
)

type Service struct {
	DB *storage.DataBase
}

func NewService(ctx context.Context) (*Service, error) {
	db, err := storage.NewDataBase(ctx)

	if err != nil {
		log.Fatalf("Failed to db connect: %v", err)
	}

	return &Service{
		DB: db,
	}, nil
}

func (s *Service) Run() {
	for {
	}
}
