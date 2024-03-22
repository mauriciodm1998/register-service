package sqs

import (
	"context"
	"encoding/json"
	"register-service/internal/channels"
	"register-service/internal/config"
	"register-service/internal/domain"
	"register-service/internal/integration/mail"
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
	sqsService     *sqs.SQS
	service        service.RegisterService
	queueProcessor map[string]func(queue string, ch chan *sqs.Message)
}

func NewSQS(mailer mail.Mailer) channels.Channel {
	once.Do(func() {
		sess := session.Must(session.NewSessionWithOptions(session.Options{
			Config: aws.Config{
				Endpoint:   aws.String(config.Get().SQS.Endpoint),
				Region:     aws.String(config.Get().SQS.Region),
				DisableSSL: aws.Bool(true),
			},
		}))

		sqs := &queueSQS{
			sqsService:     sqs.New(sess),
			service:        service.NewRegisterService(mailer),
			queueProcessor: make(map[string]func(queue string, ch chan *sqs.Message)),
		}

		sqs.queueProcessor[config.Get().SQS.ClockInQueue] = sqs.processClockInMessage
		sqs.queueProcessor[config.Get().SQS.ReportQueue] = sqs.processReportMessage

		instance = sqs
	})

	return instance
}

func (q *queueSQS) Start() error {
	for queue, processor := range q.queueProcessor {
		channel := make(chan *sqs.Message)
		go q.receiveMessage(queue, channel)
		go processor(queue, channel)
	}
	return nil
}

func (q *queueSQS) receiveMessage(queueToListen string, ch chan<- *sqs.Message) {
	for {
		paramsOrder := &sqs.ReceiveMessageInput{
			QueueUrl:            &queueToListen,
			MaxNumberOfMessages: aws.Int64(1),
		}

		resp, err := q.sqsService.ReceiveMessage(paramsOrder)
		if err != nil {
			log.Err(err).Msg("an error occurred when receive message from the queue")
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

func (q *queueSQS) processClockInMessage(queue string, ch chan *sqs.Message) {
	for {
		msg := <-ch
		var clockIn domain.ClockInRegister

		err := json.Unmarshal([]byte(*msg.Body), &clockIn)
		if err != nil {
			log.Err(err).Any("msg_id", msg.MessageId).Msg("an error occurred when reading message")
			continue
		}

		err = q.service.QueueToDatabase(context.Background(), clockIn)
		if err != nil {
			log.Err(err).Any("msg_id", msg.MessageId).Msg("an error occurred when processing message")
			continue
		}
		q.deleteMessage(msg, queue)
	}
}

func (q *queueSQS) processReportMessage(queue string, ch chan *sqs.Message) {
	for {
		msg := <-ch
		var report domain.MonthReportRequest

		err := json.Unmarshal([]byte(*msg.Body), &report)
		if err != nil {
			log.Err(err).Any("msg_id", msg.MessageId).Msg("an error occurred when reading message")
			continue
		}

		err = q.service.ReportAppointments(context.Background(), report)
		if err != nil {
			log.Err(err).Any("msg_id", msg.MessageId).Msg("an error occurred when processing message")
			continue
		}
		q.deleteMessage(msg, queue)
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
