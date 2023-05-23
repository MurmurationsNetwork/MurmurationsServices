package config

import (
	"log"
	"time"

	env "github.com/caarlos0/env/v6"
)

var Conf = config{}

type config struct {
	Server serverConf
	Mongo  mongoConf
	Static staticConf
}

type serverConf struct {
	Port                string        `env:"SERVER_PORT,required"`
	TimeoutRead         time.Duration `env:"SERVER_TIMEOUT_READ,required"`
	TimeoutWrite        time.Duration `env:"SERVER_TIMEOUT_WRITE,required"`
	TimeoutIdle         time.Duration `env:"SERVER_TIMEOUT_IDLE,required"`
	GetRateLimitPeriod  string        `env:"GET_RATE_LIMIT_PERIOD,required"`
	PostRateLimitPeriod string        `env:"POST_RATE_LIMIT_PERIOD,required"`
}

type mongoConf struct {
	USERNAME string `env:"MONGO_USERNAME,required"`
	PASSWORD string `env:"MONGO_PASSWORD,required"`
	HOST     string `env:"MONGO_HOST,required"`
	DBName   string `env:"MONGO_DB_NAME,required"`
}

type staticConf struct {
	StaticFilePath string `env:"STATIC_FILE_PATH,required"`
}

func Init() {
	err := env.Parse(&Conf)
	if err != nil {
		log.Fatalf("Failed to decode environment variables: %s", err)
	}
}
