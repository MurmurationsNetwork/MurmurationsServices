package redis

import (
	"time"
)

type redismock struct {
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
