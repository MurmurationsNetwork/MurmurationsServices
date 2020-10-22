package nats

import (
	"fmt"
	"time"

	"github.com/MurmurationsNetwork/MurmurationsServices/common/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/randstr_util"
	"github.com/cenkalti/backoff"
	"github.com/nats-io/stan.go"
)

var (
	client stan.Conn
)

func init() {
	opts := []stan.Option{stan.NatsURL("http://nats-svc:4222")}

	op := func() error {
		var err error
		client, err = stan.Connect("murmurations", randstr_util.String(8), opts...)
		if err != nil {
			return err
		}
		return nil
	}
	notify := func(err error, time time.Duration) {
		fmt.Printf("trying to re-connect NATS %s \n", err)
		fmt.Printf("retry in %s \n", time)
	}

	b := backoff.NewExponentialBackOff()
	b.MaxElapsedTime = 2 * time.Minute
	err := backoff.RetryNotify(op, b, notify)
	if err != nil {
		logger.Panic("error when trying to connect nats", err)
	}
}

func Client() stan.Conn {
	return client
}

func Disconnect() {
	logger.Info("trying to disconnect from NATS")
	client.Close()
}
