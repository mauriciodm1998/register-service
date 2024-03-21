package domain

import "time"

type ClockInRegister struct {
	Id        int       `bson:"_id"`
	UserId    int       `bson:"user_id"`
	Date      time.Time `bson:"date"`
	Time      time.Time `bson:"time"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}
