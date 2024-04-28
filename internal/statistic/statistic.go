package statistic

import (
	"Service-oriented-architectures/internal/common"
	"Service-oriented-architectures/internal/statistic/storage"
	"context"

	"encoding/json"
	"log"

	"github.com/IBM/sarama"
)

type Service struct {
	DB                *storage.DataBase
	StatisticConsumer sarama.Consumer
	UpdateConsumer    sarama.PartitionConsumer
}

func NewService(ctx context.Context) (*Service, error) {
	consumer, err := sarama.NewConsumer([]string{"kafka:9092"}, nil)
	if err != nil {
		log.Fatalf("Failed to create consumer: %v", err)
	}

	updateConsumer, err := consumer.ConsumePartition("update", 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Failed to consume partition: %v", err)
	}

	log.Println("Consumer started")

	db, err := storage.NewDataBase(ctx)

	if err != nil {
		log.Fatalf("Failed to db connect: %v", err)
	}

	return &Service{
		DB:                db,
		StatisticConsumer: consumer,
		UpdateConsumer:    updateConsumer,
	}, nil
}

func (s *Service) Run(ctx context.Context) {
	for {
		select {
		case msg, ok := <-s.UpdateConsumer.Messages():
			if !ok {
				log.Println("Channel closed, exiting goroutine")
				return
			}
			responseID := string(msg.Key)
			log.Printf("Received message with id %s\n", responseID)

			var newUpdate common.PostStatistic
			err := json.Unmarshal(msg.Value, &newUpdate)
			if err != nil {
				log.Printf("Failed to unmarshal update: %v", err)
				continue
			}

			if newUpdate.Operation == common.Like {
				err = s.DB.SetLike(ctx, newUpdate.PostID, newUpdate.UserID)
				if err != nil {
					log.Printf("Failed to set like: %v", err)
				}
			} else if newUpdate.Operation == common.View {
				err = s.DB.SetView(ctx, newUpdate.PostID, newUpdate.UserID)
				if err != nil {
					log.Printf("Failed to set view: %v", err)
				}
			}
		}
	}
}
