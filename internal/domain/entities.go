package domain

import "time"

type ClockInRegister struct {
	Id        int
	UserId    int
	Date      string
	Time      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
