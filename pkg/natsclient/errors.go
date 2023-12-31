package natsclient

import (
	"errors"
)

var (
	ErrConnectReqTimeout = errors.New("natsclient: connect request timeout")
)
