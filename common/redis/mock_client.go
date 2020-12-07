package redis

import (
	"time"

	"github.com/go-redis/redis/v8"
)

type redismock struct {
	client *redis.Client
}

func (r *redismock) Ping() error {
	return nil
}

func (r *redismock) Set(key string, value interface{}, expiration time.Duration) error {
	return nil
}

func (r *redismock) Get(key string) (string, error) {
	return "", nil
}
