package main

import (
	"register-service/internal/channels/rest"
	"register-service/internal/channels/sqs"
	"register-service/internal/config"
	"register-service/internal/integration/mail"

	"github.com/rs/zerolog/log"
)

func main() {
	config.ParseFromFlags()

	mailer := mail.NewMailer(config.Get().Mailer.From, config.Get().Mailer.Pwd, config.Get().Mailer.Address)
	go func() {
		log.Fatal().Err(sqs.NewSQS(mailer).Start())
	}()

	log.Fatal().Err(rest.NewRegisterChannel(mailer).Start())
}
