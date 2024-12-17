package config

import (
	"time"
)

var Values = config{}

type config struct {
	Server    serverConf
	Index     indexConf
	Library   libraryConf
	DataProxy dataProxyConf
	Mongo     mongoConf
}

type serverConf struct {
	Port                string        `env:"SERVER_PORT,required"`
	TimeoutRead         time.Duration `env:"SERVER_TIMEOUT_READ,required"`
	TimeoutWrite        time.Duration `env:"SERVER_TIMEOUT_WRITE,required"`
	TimeoutIdle         time.Duration `env:"SERVER_TIMEOUT_IDLE,required"`
	GetRateLimitPeriod  string        `env:"GET_RATE_LIMIT_PERIOD,required"`
	PostRateLimitPeriod string        `env:"POST_RATE_LIMIT_PERIOD,required"`
}

type indexConf struct {
	URL string `env:"INDEX_HOST,required"`
}

type libraryConf struct {
	InternalURL string `env:"LIBRARY_URL,required"`
}

type dataProxyConf struct {
	URL string `env:"EXTERNAL_DATA_PROXY_URL,required"`
}

type mongoConf struct {
	USERNAME string `env:"MONGO_USERNAME,required"`
	PASSWORD string `env:"MONGO_PASSWORD,required"`
	HOST     string `env:"MONGO_HOST,required"`
	DBName   string `env:"MONGO_DB_NAME,required"`
}
