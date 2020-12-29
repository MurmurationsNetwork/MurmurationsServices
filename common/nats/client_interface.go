package nats

import (
	"os"

	"github.com/MurmurationsNetwork/MurmurationsServices/common/backoff"
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

		operation := func() error {
			var err error
			client, err = stan.Connect(stanClusterID, clientID, opts...)
			if err != nil {
				return err
			}
			return nil
		}
		err := backoff.NewBackoff(operation, "Trying to re-connect NATS")
		if err != nil {
			return ErrConnectReqTimeout
		}
	}

	Client.setClient(client)

	return nil
}
