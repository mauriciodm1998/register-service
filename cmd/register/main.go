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

	now := time.Now().UTC()
	rep := repository.New()
	rep.Create(context.Background(), domain.ClockInRegister{
		Id:        int(uuid.New().ID()),
		UserId:    int(uuid.New().ID()),
		Date:      time.Now().UTC().Truncate(24 * time.Hour),
		Time:      time.Now().UTC(),
		CreatedAt: now.UTC(),
	})

	rep.GetMonthAppointments(context.Background(), 475075755, time.Now().UTC())
}
