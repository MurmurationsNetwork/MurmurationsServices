package redis

import (
	"context"
	"time"

	redis "github.com/go-redis/redis/v8"
)

type redisImpl struct {
	client *redis.Client
}

func (r *redisImpl) Ping() error {
	ping := r.client.Ping(context.Background())
	if err := ping.Err(); err != nil {
		return err
	}
	return nil
}

func (r *redisImpl) Set(key string, value interface{}, expiration time.Duration) error {
	err := r.client.Set(context.Background(), key, value, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *redisImpl) Get(key string) (string, error) {
	get := r.client.Get(context.Background(), key)
	if err := get.Err(); err != nil {
		if err == redis.Nil {
			return "", nil
		}
		return "", err
	}
	return get.Val(), nil
}
