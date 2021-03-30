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
	Port            string        `env:"SERVER_PORT,required"`
	TimeoutRead     time.Duration `env:"SERVER_TIMEOUT_READ,required"`
	TimeoutWrite    time.Duration `env:"SERVER_TIMEOUT_WRITE,required"`
	TimeoutIdle     time.Duration `env:"SERVER_TIMEOUT_IDLE,required"`
	RateLimitPeriod string        `env:"RATE_LIMIT_PERIOD,required"`
}

type mongoConf struct {
	USERNAME string `env:"MONGO_USERNAME,required"`
	PASSWORD string `env:"MONGO_PASSWORD,required"`
	HOST     string `env:"MONGO_HOST,required"`
	DBName   string `env:"MONGO_DB_NAME,required"`
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
