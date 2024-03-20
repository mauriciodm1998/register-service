package rest

import "register-service/internal/domain"

func (r RegisterRequest) ToClockInRegister() *domain.ClockInRegister {
	return &domain.ClockInRegister{
		Date:   r.Date,
		Time:   r.Time,
		UserId: 1, // TODO: extract from request token
	}
}
