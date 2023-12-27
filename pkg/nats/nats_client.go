package nats

import (
	"context"
	"fmt"
	"os"

	"github.com/nats-io/nats.go"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/event"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/retry"
)

const StreamName = "EVENTS" // Name of the stream in NATS.

var Client natsClientInterface // Global client interface.

type natsClient struct {
	conn      *nats.Conn
	jsContext nats.JetStreamContext
	consumers []*nats.ConsumerInfo
}

// NewClient initializes the NATS client.
func NewClient(natsURL string) error {
	if isTestEnvironment() {
		Client = &mockClient{}
		return nil
	}

	client := &natsClient{}
	Client = client

	conn, jsContext, err := client.setupNATSConnection(natsURL)
	if err != nil {
		return err
	}

	client.setConnection(conn, jsContext)
	return nil
}

// JetStream returns the JetStream context.
func (c *natsClient) JetStream() nats.JetStreamContext {
	return c.jsContext
}

// SubscribeToSubjects sets up subscriptions to given subjects.
func (c *natsClient) SubscribeToSubjects(subjects ...event.Subject) error {
	ctx := context.Background()
	subjectStrings := make([]string, len(subjects))
	for i, subject := range subjects {
		subjectStrings[i] = string(subject)
	}

	err := c.createStreamAndConsumers(ctx, c.jsContext, subjectStrings)
	if err != nil {
		return err
	}

	return nil
}

// Disconnect closes the NATS connection and deletes consumers.
func (c *natsClient) Disconnect() {
	for _, consumer := range c.consumers {
		_ = c.jsContext.DeleteConsumer(StreamName, consumer.Name)
	}
	if c.conn != nil {
		_ = c.conn.Drain()
		c.conn.Close()
	}
}

// isTestEnvironment checks for a test environment.
func isTestEnvironment() bool {
	return os.Getenv("APP_ENV") == "test"
}

// setupNATSConnection establishes a connection to NATS and sets up JetStream.
func (c *natsClient) setupNATSConnection(
	natsURL string,
) (*nats.Conn, nats.JetStreamContext, error) {
	conn, err := c.connectToNATS(natsURL)
	if err != nil {
		return nil, nil, err
	}

	jsContext, err := c.setupJetStream(conn)
	if err != nil {
		return nil, nil, err
	}

	return conn, jsContext, nil
}

// connectToNATS establishes a connection to a NATS server.
func (c *natsClient) connectToNATS(natsURL string) (*nats.Conn, error) {
	var conn *nats.Conn
	err := retry.Do(func() error {
		var err error
		conn, err = nats.Connect(natsURL)
		return err
	})
	if err != nil {
		return nil, fmt.Errorf(
			"failed to connect to NATS after retries: %w", err,
		)
	}
	return conn, nil
}

// setupJetStream initializes the JetStream context.
func (c *natsClient) setupJetStream(
	conn *nats.Conn,
) (nats.JetStreamContext, error) {
	jsContext, err := conn.JetStream()
	if err != nil {
		return nil, fmt.Errorf(
			"failed to initialize JetStream context: %w", err,
		)
	}
	return jsContext, nil
}

// setConnection sets the NATS connection and JetStream context.
func (c *natsClient) setConnection(
	conn *nats.Conn,
	jsContext nats.JetStreamContext,
) {
	c.conn = conn
	c.jsContext = jsContext
}

// createConsumer adds a new consumer to the JetStream.
func (c *natsClient) createConsumer(
	jsContext nats.JetStreamContext,
	streamName, subject, durableName string,
) error {
	consumerConfig := &nats.ConsumerConfig{
		Durable:        durableName,
		FilterSubject:  subject,
		AckPolicy:      nats.AckExplicitPolicy,
		DeliverSubject: subject,
	}
	consumerInfo, err := jsContext.AddConsumer(streamName, consumerConfig)

	c.consumers = append(c.consumers, consumerInfo)

	if err != nil {
		return fmt.Errorf("error creating consumer: %v", err)
	}
	return nil
}

// createStreamAndConsumers creates a stream and necessary consumers.
func (c *natsClient) createStreamAndConsumers(
	ctx context.Context,
	jsContext nats.JetStreamContext,
	subjects []string,
) error {
	if err := c.ensureStreamExists(ctx, jsContext, subjects); err != nil {
		return err
	}

	for _, subject := range subjects {
		consumerErr := c.createConsumer(
			jsContext,
			StreamName,
			subject,
			subject+"_consumer",
		)
		if consumerErr != nil {
			return fmt.Errorf(
				"failed to create consumer for subject %s: %w",
				subject, consumerErr,
			)
		}
	}

	return nil
}

// ensureStreamExists checks and creates a stream if it doesn't exist.
func (c *natsClient) ensureStreamExists(
	ctx context.Context,
	jsContext nats.JetStreamContext,
	subjects []string,
) error {
	_, err := jsContext.StreamInfo(StreamName)
	if err == nil {
		return nil // Stream already exists.
	}
	if err != nats.ErrStreamNotFound {
		return fmt.Errorf("error checking stream existence: %v", err)
	}
	return c.createStream(ctx, jsContext, subjects)
}

// createStream configures and adds a new stream to JetStream.
func (c *natsClient) createStream(
	_ context.Context,
	jsContext nats.JetStreamContext,
	subjects []string,
) error {
	streamConfig := &nats.StreamConfig{
		Name:              StreamName,
		Subjects:          subjects,
		Retention:         nats.InterestPolicy,
		Discard:           nats.DiscardOld,
		Storage:           nats.FileStorage,
		MaxMsgsPerSubject: 1000,
		MaxMsgSize:        1 << 20, // 1 MB
		NoAck:             false,
	}
	_, err := jsContext.AddStream(streamConfig)
	if err != nil {
		return fmt.Errorf("error creating stream: %v", err)
	}
	return nil
}
