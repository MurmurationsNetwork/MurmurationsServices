package config

import (
	"log"

	env "github.com/caarlos0/env/v6"
)

var Conf = config{}

type config struct {
	Mongo      mongoConf
	ES         esConf
	TTL        int `env:"TTL,required"`
	DELETEDTTL int `env:"DELETED_TTL,required"`
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

func Init() {
	err := env.Parse(&Conf)
	if err != nil {
		log.Fatalf("Failed to decode environment variables: %s", err)
	}
}
