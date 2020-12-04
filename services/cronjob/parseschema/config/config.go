package config

import (
	"log"

	"github.com/caarlos0/env/v6"
)

var Conf = config{}

type config struct {
	CDN   cdnConf
	Mongo mongoConf
}

type mongoConf struct {
	URL    string `env:"LIBRARY_MONGO_URL,required"`
	DBName string `env:"LIBRARY_MONGO_DB_NAME,required"`
}

type cdnConf struct {
	URL string `env:"LIBRARY_CDN_URL,required"`
}

func Init() {
	err := env.Parse(&Conf)
	if err != nil {
		log.Fatalf("Failed to decode environment variables: %s", err)
	}
}
