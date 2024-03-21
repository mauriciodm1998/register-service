package service

import (
	"context"
	"fmt"
	"register-service/internal/config"
	"register-service/internal/domain"
	"register-service/internal/integration/sqs_publisher"
	"time"

	"github.com/google/uuid"
)

type RegisterService interface {
	ClockIn(ctx context.Context, userId int) error
}

type registerService struct {
	publisher sqs_publisher.Publisher
}

func NewRegisterService() RegisterService {
	return &registerService{
		publisher: sqs_publisher.NewSQS(),
	}
}

func (s *registerService) ClockIn(ctx context.Context, userId int) error {
	if err := s.publisher.SendMessage(
		domain.ClockInRegister{
			Id:        uuid.New().String(),
			UserId:    userId,
			Date:      time.Now().Truncate(24 * time.Hour),
			Time:      time.Now(),
			CreatedAt: time.Now(),
		},
		config.Get().SQS.ClockInQueue,
	); err != nil {
		return fmt.Errorf("an error occurred when creating a clock in register: %w", err)
	}

	return nil
}
