package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type redisImpl struct {
	client *redis.Client
}

func (r *redisImpl) Ping() error {
	ping := r.client.Ping(context.Background())
	return ping.Err()
}

func (r *redisImpl) Set(
	key string,
	value interface{},
	expiration time.Duration,
) error {
	return r.client.Set(context.Background(), key, value, expiration).Err()
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
