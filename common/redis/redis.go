package redis

import (
	"os"
	"time"

	redis "github.com/go-redis/redis/v8"
)

type Redis interface {
	Ping() error
	Set(key string, value interface{}, expiration time.Duration) error
	Get(key string) (string, error)
}

func NewClient(url string) Redis {
	if os.Getenv("ENV") == "test" {
		return &redismock{}
	}
	return &redisImpl{
		client: redis.NewClient(&redis.Options{
			Addr:     url,
			Password: "",
			DB:       0,
		}),
	}
}
