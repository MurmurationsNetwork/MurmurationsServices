package config

import (
	"log"

	env "github.com/caarlos0/env/v10"
)

var Values = config{}

type config struct {
	Mongo mongoConf
	ES    esConf
	TTL   ttlConf
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

type ttlConf struct {
	ValidationFailedTTL int64 `env:"VALIDATION_FAILED_TTL,required"`
	DeletedTTL          int64 `env:"DELETED_TTL,required"`
}

func Init() {
	err := env.Parse(&Values)
	if err != nil {
		log.Fatalf("Failed to decode environment variables: %s", err)
	}
}
