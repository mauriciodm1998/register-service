package sqs_publisher

import (
	"encoding/json"
	"register-service/internal/config"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type queueSQS struct {
	queueSvc *sqs.SQS
}

type Publisher interface {
	SendMessage(inputMsg any, queueURL string) error
}

func NewSQS() Publisher {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region:     aws.String(config.Get().AWS.Region),
			DisableSSL: aws.Bool(true),
		},
	}))

	return &queueSQS{
		queueSvc: sqs.New(sess),
	}
}

func (q *queueSQS) SendMessage(inputMsg any, queueURL string) error {
	msg, err := json.Marshal(inputMsg)
	if err != nil {
		return err
	}

	params := &sqs.SendMessageInput{
		QueueUrl:    &queueURL,
		MessageBody: aws.String(string(msg)),
	}

	_, err = q.queueSvc.SendMessage(params)
	if err != nil {
		return err
	}

	return nil
}
