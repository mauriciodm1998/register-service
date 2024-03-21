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
	Database struct {
		ConnectionString string `cfg:"connection_string"`
	} `cfg:"database"`
	Token struct {
		Key string `cfg:"key"`
	} `cfg:"token"`
	SQS struct {
		ClockInQueue string `cfg:"clock_in_queue"`
		Region       string `cfg:"region"`
		Endpoint     string `cfg:"endpoint"`
	} `cfg:"sqs"`
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
		cfg.Dirs(dirs...),
	); err != nil {
		log.Panic().Err(err)
	}
}
