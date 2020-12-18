package config

import (
	"log"

	"github.com/caarlos0/env/v6"
)

var Conf = config{}

type config struct {
	Mongo mongoConf
	TTL   int `env:"TTL,required"`
}

type mongoConf struct {
	URL    string `env:"INDEX_MONGO_URL,required"`
	DBName string `env:"INDEX_MONGO_DB_NAME,required"`
}

func Init() {
	err := env.Parse(&Conf)
	if err != nil {
		log.Fatalf("Failed to decode environment variables: %s", err)
	}
}
