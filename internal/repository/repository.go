package repository

import (
	"context"
	"register-service/internal/config"
	"register-service/internal/domain"
	"strconv"
	"time"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/rs/zerolog/log"
)

const (
	tableName = "clock_in_register"
	region    = "us-east-1"
	index     = "user_id-date-index"
)

type Repository interface {
	Create(ctx context.Context, register domain.ClockInRegister) error
	GetDayAppointments(ctx context.Context, userId int, target time.Time) ([]domain.ClockInRegister, error)
	GetMonthAppointments(ctx context.Context, userId int, target time.Time) ([]domain.ClockInRegister, error)
	GetWeekAppointments(ctx context.Context, userId int, target time.Time) ([]domain.ClockInRegister, error)
}

type repository struct {
	database  *dynamodb.Client
	tableName string
	index     string
}

func New() Repository {
	cfg, err := awsconfig.LoadDefaultConfig(context.TODO(),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(config.Get().AWS.AccessKeyId, config.Get().AWS.SecretAccessKey, config.Get().AWS.SessionToken)), awsconfig.WithRegion(region),
	)
	if err != nil {
		log.Fatal().Err(err).Msg("an error occurred when connect to the database")
	}
	return &repository{
		database:  dynamodb.NewFromConfig(cfg),
		tableName: tableName,
		index:     index,
	}
}

func (r *repository) Create(ctx context.Context, register domain.ClockInRegister) error {
	av, err := attributevalue.MarshalMap(register)
	if err != nil {
		return err
	}

	_, err = r.database.PutItem(ctx, &dynamodb.PutItemInput{
		Item:      av,
		TableName: &r.tableName,
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) GetDayAppointments(ctx context.Context, userId int, target time.Time) ([]domain.ClockInRegister, error) {
	result, err := r.database.Query(context.TODO(), &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		IndexName:              &r.index,
		KeyConditionExpression: aws.String("user_id = :id AND #date = :date"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":id":   &types.AttributeValueMemberN{Value: strconv.Itoa(userId)},
			":date": &types.AttributeValueMemberS{Value: target.Format(time.RFC3339)},
		},
		ExpressionAttributeNames: map[string]string{
			"#date": "date",
		},
	})
	if err != nil {
		return nil, err
	}

	var appointments []domain.ClockInRegister

	for _, item := range result.Items {
		var appointment domain.ClockInRegister

		if err := attributevalue.UnmarshalMap(item, &appointment); err != nil {
			return nil, err
		}

		appointments = append(appointments, appointment)
	}

	return appointments, nil
}

func (r *repository) GetMonthAppointments(ctx context.Context, userId int, target time.Time) ([]domain.ClockInRegister, error) {
	start := time.Date(target.Year(), target.Month()-1, 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, -1)

	result, err := r.database.Query(context.TODO(), &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		IndexName:              &r.index,
		KeyConditionExpression: aws.String("user_id = :id AND #date BETWEEN :end AND :start"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":id":    &types.AttributeValueMemberN{Value: strconv.Itoa(userId)},
			":start": &types.AttributeValueMemberS{Value: start.Format(time.RFC3339)},
			":end":   &types.AttributeValueMemberS{Value: end.Format(time.RFC3339)},
		},
		ExpressionAttributeNames: map[string]string{
			"#date": "date",
		},
	})
	if err != nil {
		return nil, nil
	}

	var appointments []domain.ClockInRegister

	for _, item := range result.Items {
		var appointment domain.ClockInRegister

		if err := attributevalue.UnmarshalMap(item, &appointment); err != nil {
			return nil, err
		}

		appointments = append(appointments, appointment)
	}

	return appointments, nil
}

func (r *repository) GetWeekAppointments(ctx context.Context, userId int, target time.Time) ([]domain.ClockInRegister, error) {
	start := target.AddDate(0, 0, -int(target.Weekday()))

	result, err := r.database.Query(context.TODO(), &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		IndexName:              &r.index,
		KeyConditionExpression: aws.String("user_id = :id AND #date BETWEEN :start AND :target"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":id":     &types.AttributeValueMemberN{Value: strconv.Itoa(userId)},
			":start":  &types.AttributeValueMemberS{Value: start.Format(time.RFC3339)},
			":target": &types.AttributeValueMemberS{Value: target.Format(time.RFC3339)},
		},
		ExpressionAttributeNames: map[string]string{
			"#date": "date",
		},
	})
	if err != nil {
		return nil, nil
	}

	var appointments []domain.ClockInRegister

	for _, item := range result.Items {
		var appointment domain.ClockInRegister

		if err := attributevalue.UnmarshalMap(item, &appointment); err != nil {
			return nil, err
		}

		appointments = append(appointments, appointment)
	}

	return appointments, nil
}
