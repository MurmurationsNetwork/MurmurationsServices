package natsclient

import (
	"fmt"
	"strings"
	"sync"

	"github.com/nats-io/nats.go"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/event"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/retry"
)

// NatsClient represents the singleton NATS client.
type NatsClient struct {
	conn      *nats.Conn
	JsContext nats.JetStreamContext
	consumers []*nats.ConsumerInfo
}

var (
	instance *NatsClient
	once     sync.Once
)

// Initialize must be called at the start of your application.
func Initialize(url string) error {
	var err error
	once.Do(func() {
		instance = &NatsClient{}
		instance.conn, err = connectToNATS(url)
		if err == nil {
			instance.JsContext, err = instance.conn.JetStream()
		}
	})
	return err
}

// GetInstance returns the singleton instance of the NATS client.
func GetInstance() *NatsClient {
	if instance == nil || instance.conn == nil {
		panic("NATS client is not initialized or connection is nil. " +
			"Ensure Initialize is called correctly.")
	}
	return instance
}

// connectToNATS establishes a connection to a NATS server.
func connectToNATS(natsURL string) (*nats.Conn, error) {
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

// SubscribeToSubjects sets up subscriptions to given subjects.
func (c *NatsClient) SubscribeToSubjects(subjects ...event.Subject) error {
	subjectStrings := make([]string, len(subjects))
	for i, subject := range subjects {
		subjectStrings[i] = string(subject)
	}

	return c.createStreamAndConsumers(subjectStrings)
}

// createStreamAndConsumers creates a stream and necessary consumers.
func (c *NatsClient) createStreamAndConsumers(subjects []string) error {
	if err := c.ensureStreamExists(); err != nil {
		return err
	}

	for _, subject := range subjects {
		if err := c.createConsumer(subject, subject+"_consumer"); err != nil {
			return fmt.Errorf(
				"failed to create consumer for subject %s: %w",
				subject,
				err,
			)
		}
	}
	return nil
}

// ensureStreamExists checks and creates a stream if it doesn't exist.
func (c *NatsClient) ensureStreamExists() error {
	_, err := c.JsContext.StreamInfo(streamName)
	if err == nil {
		return nil // Stream already exists.
	}
	if err != nats.ErrStreamNotFound {
		return fmt.Errorf("error checking stream existence: %v", err)
	}
	return c.createStream()
}

// createStream configures and adds a new stream to JetStream.
func (c *NatsClient) createStream() error {
	streamConfig := &nats.StreamConfig{
		Name:              streamName,
		Subjects:          []string{"node:*"},
		Retention:         nats.WorkQueuePolicy,
		Discard:           nats.DiscardOld,
		Storage:           nats.FileStorage,
		MaxMsgsPerSubject: 1000,
		MaxMsgSize:        1 << 20, // 1 MB
		NoAck:             false,
	}
	_, err := c.JsContext.AddStream(streamConfig)
	if err != nil {
		return fmt.Errorf("error creating stream: %v", err)
	}
	return nil
}

// createConsumer adds a new consumer to the JetStream.
func (c *NatsClient) createConsumer(subject, durableName string) error {
	consumerConfig := &nats.ConsumerConfig{
		// 'Durable' names the consumer, allowing it to be durable.
		// This means the state of the consumer (like acked messages) is
		// maintained across restarts.
		Durable: durableName,
		// 'FilterSubject' specifies the subject (or subjects) this consumer
		// will listen to.
		FilterSubject: subject,
		// 'DeliverSubject' is the NATS subject where messages will be delivered
		DeliverSubject: subject,
		// 'AckExplicitPolicy' requires explicit acknowledgment of each message.
		AckPolicy: nats.AckExplicitPolicy,
	}
	consumerInfo, err := c.JsContext.AddConsumer(streamName, consumerConfig)

	if err != nil {
		return fmt.Errorf("error creating consumer: %v", err)
	}

	c.consumers = append(c.consumers, consumerInfo)
	return nil
}

// Disconnect closes the NATS connection and deletes consumers.
func (c *NatsClient) Disconnect() error {
	var errStrings []string

	// Delete consumers and collect errors.
	for _, consumer := range c.consumers {
		if err := c.JsContext.DeleteConsumer(streamName, consumer.Name); err != nil {
			errStrings = append(
				errStrings,
				fmt.Sprintf(
					"error deleting consumer %s: %v",
					consumer.Name,
					err,
				),
			)
		}
	}

	// Attempt to drain and close the connection.
	if c.conn != nil {
		if err := c.conn.Drain(); err != nil {
			errStrings = append(
				errStrings,
				fmt.Sprintf("error draining connection: %v", err),
			)
		}
		c.conn.Close()
	}

	// If there were any errors, return a combined error.
	if len(errStrings) > 0 {
		return fmt.Errorf(
			"disconnect encountered issues: %s",
			strings.Join(errStrings, ", "),
		)
	}

	return nil
}
