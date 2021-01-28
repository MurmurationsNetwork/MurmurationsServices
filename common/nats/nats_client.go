package nats

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/common/logger"
	"github.com/nats-io/stan.go"
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
