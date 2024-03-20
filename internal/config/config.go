package config

import (
	"flag"
	"log"

	"github.com/notnull-co/cfg"
)

var (
	configuration config
)

type config struct {
	Database struct {
		ConnectionString string `cfg:"connection_string"`
	} `cfg:"database"`
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
		log.Panic(err)
	}
}
