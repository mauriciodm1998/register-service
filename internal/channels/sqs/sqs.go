package sqs

import (
	"context"
	"encoding/json"
	"register-service/internal/channels"
	"register-service/internal/config"
	"register-service/internal/domain"
	"register-service/internal/service"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/rs/zerolog/log"
)

const (
	PAYMENT = "payment"
	ORDER   = "order"
)

type queueSQS struct {
	sqsService    *sqs.SQS
	service       service.RegisterService
	queuesAddress string
}

func NewSQS() channels.Channel {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Endpoint:   aws.String(config.Get().SQS.Endpoint),
			Region:     aws.String(config.Get().SQS.Region),
			DisableSSL: aws.Bool(true),
		},
	}))

	return &queueSQS{
		sqsService:    sqs.New(sess),
		service:       service.NewRegisterService(),
		queuesAddress: config.Get().SQS.ClockInQueue,
	}
}

func (q *queueSQS) Start() error {
	q.ReceiveMessage()
	return nil
}

func (q *queueSQS) ReceiveMessage() {
	for {
		paramsOrder := &sqs.ReceiveMessageInput{
			QueueUrl:            &q.queuesAddress,
			MaxNumberOfMessages: aws.Int64(1),
		}

		resp, err := q.sqsService.ReceiveMessage(paramsOrder)
		if err != nil {
			log.Err(err).Msg("an error occurred when receive message from the queue")
			continue
		}

		if len(resp.Messages) > 0 {
			for _, msg := range resp.Messages {

				err := q.processMessage([]byte(*msg.Body))
				if err != nil {
					log.Err(err).Any("msg_id", msg.MessageId).Msg("an error occurred when process message")
					continue
				}

				_, err = q.sqsService.DeleteMessage(&sqs.DeleteMessageInput{
					QueueUrl:      &q.queuesAddress,
					ReceiptHandle: msg.ReceiptHandle,
				})
				if err != nil {
					log.Err(err).Any("msg_id", msg.MessageId).Msg("an error occurred when delete message")
					continue
				}
			}
		} else {
			log.Info().Msg("no new messages")
			time.Sleep(time.Second * 10)
		}
	}
}

func (q *queueSQS) processMessage(msg []byte) error {
	var clockIn domain.ClockInRegister

	err := json.Unmarshal(msg, &clockIn)
	if err != nil {
		return err
	}

	err = q.service.QueueToDatabase(context.Background(), clockIn)
	if err != nil {
		return err
	}

	return nil
}
