package nats

import (
	stan "github.com/nats-io/stan.go"

	"github.com/MurmurationsNetwork/MurmurationsServices/common/logger"
)

type natsClient struct {
	client stan.Conn
}

func (c *natsClient) Client() stan.Conn {
	return c.client
}

func (c *natsClient) Disconnect() {
	err := c.client.Close()
	if err != nil {
		logger.Error("Error when trying to disconnect from Nats", err)
	}
}

func (c *natsClient) setClient(client stan.Conn) {
	c.client = client
}
