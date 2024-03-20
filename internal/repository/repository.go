package repository

import (
	"context"
	"register-service/internal/config"
	"register-service/internal/domain"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	collection = "order"
	database   = "order"
)

type Repository interface {
	Create(ctx context.Context, register domain.ClockInRegister) error
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

	return nil
}
