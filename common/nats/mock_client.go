package nats

import (
	stan "github.com/nats-io/stan.go"
)

type mockClient struct {
	client stan.Conn
}

func (c *mockClient) Client() stan.Conn {
	return c.client
}

func (c *mockClient) Disconnect() {
}

func (c *mockClient) setClient(client stan.Conn) {
	c.client = client
}
