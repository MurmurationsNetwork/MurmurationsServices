package nats

import (
	"errors"
)

var (
	ErrConnectReqTimeout = errors.New("stan: connect request timeout")
)
