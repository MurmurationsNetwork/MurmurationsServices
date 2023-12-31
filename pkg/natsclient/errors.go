package natsclient

import (
	"errors"
)

var (
	ErrConnectReqTimeout = errors.New("stan: connect request timeout")
)
