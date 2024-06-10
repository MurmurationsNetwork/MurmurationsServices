package redis

import (
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

type Redis interface {
	Ping() error
	Set(key string, value interface{}, expiration time.Duration) error
	Get(key string) (string, error)
}

func NewClient(url string) Redis {
	if os.Getenv("APP_ENV") == "test" {
		return &redismock{}
	}
	return &redisImpl{
		client: redis.NewClient(&redis.Options{
			Addr:         url,
			Password:     "",
			DB:           0,
			DialTimeout:  10 * time.Second,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		}),
	}
}
