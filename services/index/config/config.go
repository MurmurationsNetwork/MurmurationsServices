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
	Nats   natsConf
}

type serverConf struct {
	Port         string        `env:"SERVER_PORT,required"`
	TimeoutRead  time.Duration `env:"SERVER_TIMEOUT_READ,required"`
	TimeoutWrite time.Duration `env:"SERVER_TIMEOUT_WRITE,required"`
	TimeoutIdle  time.Duration `env:"SERVER_TIMEOUT_IDLE,required"`
}

type mongoConf struct {
	URL    string `env:"MONGO_URL,required"`
	DBName string `env:"MONGO_DB_NAME,required"`
}

type esConf struct {
	URL string `env:"ELASTICSEARCH_URL,required"`
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
