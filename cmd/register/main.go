package main

import (
	"register-service/internal/channels/rest"
	"register-service/internal/channels/sqs"
	"register-service/internal/config"

	"github.com/rs/zerolog/log"
)

func main() {
	config.ParseFromFlags()

	go func() {
		log.Fatal().Err(sqs.NewSQS().Start())
	}()

	log.Fatal().Err(rest.NewRegisterChannel().Start())
}
