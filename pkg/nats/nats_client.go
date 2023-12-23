package nats

import (
	"github.com/nats-io/nats.go"
)

type natsClient struct {
	conn *nats.Conn
	js   nats.JetStreamContext
}

// Client returns the JetStream context.
func (c *natsClient) JetStream() nats.JetStreamContext {
	return c.js
}

// Disconnect closes the NATS connection.
func (c *natsClient) Disconnect() {
	if c.conn != nil {
		c.conn.Close()
	}
}

// setClient sets the NATS connection.
func (c *natsClient) setClient(conn *nats.Conn, js nats.JetStreamContext) {
	c.conn = conn
	c.js = js
}
