package redis

import (
	"time"
)

type redismock struct {
}

func (*redismock) Ping() error {
	return nil
}

func (*redismock) Set(
	_ string,
	_ interface{},
	_ time.Duration,
) error {
	return nil
}

func (*redismock) Get(_ string) (string, error) {
	return "", nil
}
