package nats

import (
	"fmt"
	"os"
	"time"

	"github.com/MurmurationsNetwork/MurmurationsServices/common/logger"
	"github.com/cenkalti/backoff"
	"github.com/nats-io/stan.go"
)

var (
	Client natsClientInterface
)

type natsClientInterface interface {
	Client() stan.Conn
	Disconnect()

	setClient(stan.Conn)
}

func init() {
	if os.Getenv("ENV") == "test" {
		Client = &mockClient{}
		return
	}
	Client = &natsClient{}
}

func NewClient(stanClusterID, clientID, natsURL string) error {
	var client stan.Conn

	if os.Getenv("ENV") != "test" {
		opts := []stan.Option{stan.NatsURL(natsURL)}

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
		}

		b := backoff.NewExponentialBackOff()
		b.MaxElapsedTime = 2 * time.Minute
		err := backoff.RetryNotify(op, b, notify)
		if err != nil {
			return ErrConnectReqTimeout
		}

		fmt.Println(client)
	}

	Client.setClient(client)

	return nil
}
