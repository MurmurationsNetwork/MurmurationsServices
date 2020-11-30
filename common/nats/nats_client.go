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
	logger.Info("trying to disconnect from NATS")
	c.client.Close()
}

func (c *natsClient) setClient(client stan.Conn) {
	c.client = client
}
