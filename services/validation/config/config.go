package config

import (
	"log"
	"time"

	"github.com/caarlos0/env/v6"
)

var Conf = config{}

type config struct {
	Server  serverConf
	Library libraryConf
	Nats    natsConf
}

type serverConf struct {
	Port         string        `env:"SERVER_PORT,required"`
	TimeoutRead  time.Duration `env:"SERVER_TIMEOUT_READ,required"`
	TimeoutWrite time.Duration `env:"SERVER_TIMEOUT_WRITE,required"`
	TimeoutIdle  time.Duration `env:"SERVER_TIMEOUT_IDLE,required"`
}

type libraryConf struct {
	URL         string `env:"LIBRARY_CDN_URL,required"`
	InternalURL string `env:"LIBRARY_URL,required"`
}

type natsConf struct {
	ClusterID string `env:"NATS_CLUSTER_ID,required"`
	ClientID  string `env:"NATS_CLIENT_ID,required"`
	URL       string `env:"NATS_URL,required"`
}

func Init() {
	err := env.Parse(&Conf)
	if err != nil {
		log.Fatalf("Failed to decode environment variables: %s", err)
	}
}
