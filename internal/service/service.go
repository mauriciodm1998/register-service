package service

import (
	"context"
)

type RegisterService interface {
	ClockIn(ctx context.Context, userId string) error
}

func NewRegisterService() RegisterService {
	return nil
}
