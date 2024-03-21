package rest

import "register-service/internal/domain"

func (r RegisterRequest) ToClockInRegister() *domain.ClockInRegister {
	return &domain.ClockInRegister{
		UserId: 1, // TODO: extract from request token
	}
}
