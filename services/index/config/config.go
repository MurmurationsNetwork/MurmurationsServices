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
	Mongo   mongoConf
	ES      esConf
	Nats    natsConf
	TTL     ttlConf
}

type serverConf struct {
	Port                string        `env:"SERVER_PORT,required"`
	TimeoutRead         time.Duration `env:"SERVER_TIMEOUT_READ,required"`
	TimeoutWrite        time.Duration `env:"SERVER_TIMEOUT_WRITE,required"`
	TimeoutIdle         time.Duration `env:"SERVER_TIMEOUT_IDLE,required"`
	GetRateLimitPeriod  string        `env:"GET_RATE_LIMIT_PERIOD,required"`
	PostRateLimitPeriod string        `env:"POST_RATE_LIMIT_PERIOD,required"`
	TagsArraySize       string        `env:"TAGS_ARRAY_SIZE,required"`
	TagsStringLength    string        `env:"TAGS_STRING_LENGTH,required"`
	TagsFuzziness       string        `env:"TAGS_FUZZINESS,required"`
}

type libraryConf struct {
	URL         string `env:"LIBRARY_CDN_URL,required"`
	InternalURL string `env:"LIBRARY_URL,required"`
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

type ttlConf struct {
	DeletedTTL int64 `env:"DELETED_TTL,required"`
}

func Init() {
	err := env.Parse(&Conf)
	if err != nil {
		log.Fatalf("Failed to decode environment variables: %s", err)
	}
}
