package sqs

import (
	"context"
	"encoding/json"
	"register-service/internal/channels"
	"register-service/internal/config"
	"register-service/internal/domain"
	"register-service/internal/service"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/rs/zerolog/log"
)

var (
	once     sync.Once
	instance channels.Channel
)

type queueSQS struct {
	sqsService *sqs.SQS
	service    service.RegisterService
}

func NewSQS() channels.Channel {
	once.Do(func() {
		sess := session.Must(session.NewSessionWithOptions(session.Options{
			Config: aws.Config{
				Region:     aws.String(config.Get().AWS.Region),
				DisableSSL: aws.Bool(true),
			},
		}))

		sqs := &queueSQS{
			sqsService: sqs.New(sess),
			service:    service.NewRegisterService(),
		}

		instance = sqs
	})

	return instance
}

func (q *queueSQS) Start() error {
	clockinChannel := make(chan *sqs.Message)
	reportChannel := make(chan *sqs.Message)

	go q.receiveMessage(config.Get().SQS.ClockInQueue, clockinChannel)
	go q.receiveMessage(config.Get().SQS.ReportQueue, reportChannel)

	q.messageProcessor(clockinChannel, reportChannel)

	return nil
}

func (q *queueSQS) messageProcessor(clockinChannel chan *sqs.Message, reportChannel chan *sqs.Message) {
	for {
		select {
		case clockMessage := <-clockinChannel:
			log.Info().Any("msg_id", clockMessage.MessageId).Msg("msg received from clock in queue")

			var clockIn domain.ClockInRegister

			err := json.Unmarshal([]byte(*clockMessage.Body), &clockIn)
			if err != nil {
				log.Err(err).Any("msg_id", clockMessage.MessageId).Msg("an error occurred when reading message")
				continue
			}

			err = q.service.QueueToDatabase(context.Background(), clockIn)
			if err != nil {
				log.Err(err).Any("msg_id", clockMessage.MessageId).Msg("an error occurred when processing message")
				continue
			}

			q.deleteMessage(clockMessage, config.Get().SQS.ClockInQueue)

		case reportMessage := <-reportChannel:
			log.Info().Any("msg_id", reportMessage.MessageId).Msg("msg received from report queue")

			var report domain.MonthReportRequest

			err := json.Unmarshal([]byte(*reportMessage.Body), &report)
			if err != nil {
				log.Err(err).Any("msg_id", reportMessage.MessageId).Msg("an error occurred when reading message")
				continue
			}

			err = q.service.ReportAppointments(context.Background(), report)
			if err != nil {
				log.Err(err).Any("msg_id", reportMessage.MessageId).Msg("an error occurred when processing message")
				continue
			}

			q.deleteMessage(reportMessage, config.Get().SQS.ReportQueue)
		}
	}
}

func (q *queueSQS) receiveMessage(queueToListen string, ch chan<- *sqs.Message) {
	for {
		paramsOrder := &sqs.ReceiveMessageInput{
			QueueUrl:            &queueToListen,
			MaxNumberOfMessages: aws.Int64(1),
		}

		resp, err := q.sqsService.ReceiveMessage(paramsOrder)
		if err != nil {
			log.Fatal().Err(err).Msg("an error occurred when receive message from the queue")
			continue
		}

		if len(resp.Messages) > 0 {
			for _, msg := range resp.Messages {
				ch <- msg
			}
		} else {
			log.Info().Msg("no new messages")
			time.Sleep(time.Second * 10)
		}
	}
}

func (q *queueSQS) deleteMessage(msg *sqs.Message, queue string) {
	_, err := q.sqsService.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      &queue,
		ReceiptHandle: msg.ReceiptHandle,
	})
	if err != nil {
		log.Err(err).Msg("an error occurred when deleting message")
	}
}
