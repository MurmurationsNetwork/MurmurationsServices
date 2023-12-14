package nats

import (
	"os"

	stan "github.com/nats-io/stan.go"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/retry"
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
	if os.Getenv("APP_ENV") == "test" {
		Client = &mockClient{}
		return
	}
	Client = &natsClient{}
}

func NewClient(stanClusterID, clientID, natsURL string) error {
	var client stan.Conn

	if os.Getenv("APP_ENV") != "test" {
		opts := []stan.Option{stan.NatsURL(natsURL)}

		operation := func() error {
			var err error
			client, err = stan.Connect(stanClusterID, clientID, opts...)
			if err != nil {
				return err
			}
			return nil
		}
		err := retry.Do(operation)
		if err != nil {
			return ErrConnectReqTimeout
		}
	}

	Client.setClient(client)

	return nil
}
