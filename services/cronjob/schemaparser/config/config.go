package config

import (
	"log"

	env "github.com/caarlos0/env/v6"
)

var Conf = config{}

type config struct {
	Library libraryConf
	Mongo   mongoConf
	Redis   redisConf
	Github  githubConf
}

type libraryConf struct {
	URL string `env:"LIBRARY_URL,required"`
}

type mongoConf struct {
	USERNAME string `env:"MONGO_USERNAME,required"`
	PASSWORD string `env:"MONGO_PASSWORD,required"`
	HOST     string `env:"MONGO_HOST,required"`
	DBName   string `env:"MONGO_DB_NAME,required"`
}

type redisConf struct {
	URL string `env:"REDIS_URL,required"`
}

type githubConf struct {
	TOKEN     string `env:"GITHUB_TOKEN,required"`
	BranchURL string `env:"GITHUB_BRANCH_URL,required"`
	TreeURL   string `env:"GITHUB_TREE_URL,required"`
}

func Init() {
	err := env.Parse(&Conf)
	if err != nil {
		log.Fatalf("Failed to decode environment variables: %s", err)
	}
}
