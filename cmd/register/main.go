package main

import (
	"context"
	"register-service/internal/config"
	"register-service/internal/domain"
	"register-service/internal/repository"
	"time"

	"github.com/google/uuid"
)

func main() {
	config.ParseFromFlags()

	now := time.Now()
	repository.New().Create(context.Background(), domain.ClockInRegister{
		Id:        int(uuid.New().ID()),
		UserId:    int(uuid.New().ID()),
		Date:      now.Format("2006-01-02"),
		Time:      now.Format("15:04"),
		CreatedAt: now,
	})
}
