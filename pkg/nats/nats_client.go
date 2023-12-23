package nats

import (
	"context"

	"github.com/nats-io/nats.go"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/event"
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

func (c *natsClient) SubscribeToSubjects(subjects ...event.Subject) error {
	ctx := context.Background()

	subjectStrings := make([]string, len(subjects))
	for i, subject := range subjects {
		subjectStrings[i] = string(subject)
	}

	if err := createStreamAndConsumers(ctx, c.js, subjectStrings); err != nil {
		return err
	}
	// Additional logic to handle subscriptions can be added here
	return nil
}

// setClient sets the NATS connection.
func (c *natsClient) setClient(conn *nats.Conn, js nats.JetStreamContext) {
	c.conn = conn
	c.js = js
}
