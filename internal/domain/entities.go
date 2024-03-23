package domain

import "time"

type ClockInRegister struct {
	Id        string    `dynamodbav:"id"`
	UserId    int       `dynamodbav:"user_id"`
	Date      time.Time `dynamodbav:"date"`
	Time      time.Time `dynamodbav:"time"`
	CreatedAt time.Time `dynamodbav:"created_at"`
	UpdatedAt time.Time `dynamodbav:"updated_at"`
}

type DailyRegister struct {
	Clocks []ClockInRegister
	Hours  int
}

type MonthReportRequest struct {
	UserId int
	Time   time.Time
	Mail   string
}
