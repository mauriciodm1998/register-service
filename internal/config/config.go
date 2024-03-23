package config

import (
	"flag"

	"github.com/notnull-co/cfg"
	"github.com/rs/zerolog/log"
)

var (
	configuration config
)

type config struct {
	Server struct {
		Port string `cfg:"port"`
	} `cfg:"server"`
	Token struct {
		Key string `cfg:"key"`
	} `cfg:"token"`
	Mailer struct {
		From string `cfg:"from"`
		Pwd  string `cfg:"pwd"`
	} `cfg:"mailer"`
	SQS struct {
		ClockInQueue string `cfg:"clock_in_queue"`
		ReportQueue  string `cfg:"report_queue"`
	} `cfg:"sqs"`
	AWS struct {
		Region          string `cfg:"region"`
		AccessKeyId     string `cfg:"access_key_id"`
		SecretAccessKey string `cfg:"secret_access_key"`
		SessionToken    string `cfg:"session_token"`
	}
}

func Get() config {
	return configuration
}

func ParseFromFlags() {
	var configDir string

	flag.StringVar(&configDir, "config-dir", "./config/", "Configuration file directory")
	flag.Parse()

	parse(configDir)
}

func parse(dirs ...string) {
	if err := cfg.Load(&configuration,
		cfg.UseEnv("APP"),
		cfg.Dirs(dirs...),
	); err != nil {
		log.Panic().Err(err)
	}
}
