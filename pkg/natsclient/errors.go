package natsclient

import (
	"errors"
)

var (
	ErrConnectReqTimeout = errors.New("nats: connect request timeout")
)
