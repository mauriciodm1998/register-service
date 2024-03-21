package repository

import (
	"context"
	"fmt"
	"register-service/internal/config"
	"register-service/internal/domain"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	collection = "clock_in_register"
	database   = "default"
)

type Repository interface {
	Create(ctx context.Context, register domain.ClockInRegister) error
	GetDayAppointments(ctx context.Context, userId int, target time.Time) ([]domain.ClockInRegister, error)
	GetMonthAppointments(ctx context.Context, userId int, target time.Time) ([]domain.ClockInRegister, error)
}

type repository struct {
	collection *mongo.Collection
}

func New() Repository {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(config.Get().Database.ConnectionString))
	if err != nil {
		panic(err)
	}

	return &repository{
		collection: client.Database(database).Collection(collection),
	}
}

func (r *repository) Create(ctx context.Context, register domain.ClockInRegister) error {
	register.Time = register.Time.Truncate(60 * time.Minute)

	_, err := r.collection.InsertOne(ctx, register)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) GetMonthAppointments(ctx context.Context, userId int, target time.Time) ([]domain.ClockInRegister, error) {
	start := time.Date(target.Year(), target.Month()-1, 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, -1)

	fmt.Println(start.String())
	fmt.Println(end.String())
	filter := bson.M{
		"$and": []bson.M{
			{
				"date": bson.M{
					"$gte": start,
				},
			},
			{
				"date": bson.M{
					"$lte": end,
				},
			},
			{
				"user_id": userId,
			},
		},
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var registers []domain.ClockInRegister

	if err = cursor.All(ctx, &registers); err != nil {
		return nil, err
	}

	return registers, nil
}

func (r *repository) GetDayAppointments(ctx context.Context, userId int, target time.Time) ([]domain.ClockInRegister, error) {
	filter := bson.M{
		"$and": []bson.M{
			{
				"date": target,
			},
			{
				"user_id": userId,
			},
		},
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var registers []domain.ClockInRegister

	if err = cursor.All(ctx, &registers); err != nil {
		return nil, err
	}

	return registers, nil
}
