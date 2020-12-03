package config

import (
	"log"
	"time"

	"github.com/caarlos0/env/v6"
)

var Conf = config{}

type config struct {
	Server serverConf
	Mongo  mongoConf
	ES     esConf
}

type serverConf struct {
	Port         string        `env:"SERVER_PORT,required"`
	TimeoutRead  time.Duration `env:"SERVER_TIMEOUT_READ,required"`
	TimeoutWrite time.Duration `env:"SERVER_TIMEOUT_WRITE,required"`
	TimeoutIdle  time.Duration `env:"SERVER_TIMEOUT_IDLE,required"`
}

type mongoConf struct {
}

type esConf struct {
}

func Init() {
	err := env.Parse(&Conf)
	if err != nil {
		log.Fatalf("Failed to decode environment variables: %s", err)
	}
}
