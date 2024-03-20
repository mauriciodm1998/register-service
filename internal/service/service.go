package service

import (
	"context"
	"register-service/internal/domain"
)

type RegisterService interface {
	ClockIn(ctx context.Context, status *domain.ClockInRegister) error
}

func NewRegisterService() RegisterService {
	return nil
}
