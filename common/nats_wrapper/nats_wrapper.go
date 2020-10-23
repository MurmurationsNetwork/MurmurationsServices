package nats_wrapper

import (
	"errors"
	"fmt"
	"time"

	"github.com/MurmurationsNetwork/MurmurationsServices/common/logger"
	"github.com/cenkalti/backoff"
	"github.com/nats-io/stan.go"
)

var (
	ErrConnectReqTimeout = errors.New("stan: connect request timeout")
)

func Connect(stanClusterID, clientID, natsURL string) (stan.Conn, error) {
	opts := []stan.Option{stan.NatsURL(natsURL)}

	var client stan.Conn

	op := func() error {
		var err error
		client, err = stan.Connect(stanClusterID, clientID, opts...)
		if err != nil {
			return err
		}
		return nil
	}
	notify := func(err error, time time.Duration) {
		logger.Error("trying to re-connect NATS %s \n", err)
		logger.Info(fmt.Sprintf("trying to re-connect NATS %s \n", err))
	}

	b := backoff.NewExponentialBackOff()
	b.MaxElapsedTime = 2 * time.Minute
	err := backoff.RetryNotify(op, b, notify)
	if err != nil {
		return nil, ErrConnectReqTimeout
	}

	return client, nil
}
