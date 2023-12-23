package nats

import (
	"context"
	"fmt"
	"os"

	"github.com/nats-io/nats.go"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/retry"
)

const (
	StreamName = "EVENTS" // Name of the stream used in NATS
)

// Global client interface variable.
var Client natsClientInterface

// init initializes the global client and sets up the required subjects.
func init() {
	Client = &natsClient{}
}

// NewClient initializes the NATS client.
func NewClient(natsURL string) error {
	if isTestEnvironment() {
		Client = &mockClient{}
		return nil
	}

	nc, js, err := setupNATSConnection(natsURL)
	if err != nil {
		return err
	}

	Client.setClient(nc, js)
	return nil
}

// isTestEnvironment checks if the application is running in a test environment.
func isTestEnvironment() bool {
	return os.Getenv("APP_ENV") == "test"
}

// setupNATSConnection handles the connection to NATS and JetStream setup.
func setupNATSConnection(
	natsURL string,
) (*nats.Conn, nats.JetStreamContext, error) {
	nc, err := connectToNATS(natsURL)
	if err != nil {
		return nil, nil, err
	}

	js, err := setupJetStream(nc)
	if err != nil {
		return nil, nil, err
	}

	return nc, js, nil
}

// connectToNATS establishes a connection to a NATS server with retries.
func connectToNATS(natsURL string) (*nats.Conn, error) {
	var nc *nats.Conn
	err := retry.Do(func() error {
		var err error
		nc, err = nats.Connect(natsURL)
		return err
	})
	if err != nil {
		return nil, fmt.Errorf(
			"failed to connect to NATS after retries: %w",
			err,
		)
	}
	return nc, nil
}

// setupJetStream initializes the JetStream context on the NATS connection.
func setupJetStream(nc *nats.Conn) (nats.JetStreamContext, error) {
	js, err := nc.JetStream()
	if err != nil {
		return nil, fmt.Errorf(
			"failed to initialize JetStream context: %w",
			err,
		)
	}
	return js, nil
}

// createStreamAndConsumers creates a stream and necessary consumers in JetStream.
func createStreamAndConsumers(
	ctx context.Context,
	js nats.JetStreamContext,
	subjects []string,
) error {
	if err := ensureStreamExists(ctx, js, subjects); err != nil {
		return err
	}

	for _, subject := range subjects {
		if err := createConsumer(js, StreamName, subject, subject+"_consumer"); err != nil {
			return fmt.Errorf(
				"failed to create consumer for subject %s: %w",
				subject,
				err,
			)
		}
	}
	return nil
}

// ensureStreamExists checks if a stream exists in JetStream, and creates it if not.
func ensureStreamExists(
	ctx context.Context,
	js nats.JetStreamContext,
	subjects []string,
) error {
	_, err := js.StreamInfo(StreamName)
	if err == nil {
		return nil // Stream already exists.
	}
	if err != nats.ErrStreamNotFound {
		return fmt.Errorf("error checking stream existence: %v", err)
	}
	return createStream(ctx, js, subjects)
}

// createStream configures and adds a new stream to JetStream.
func createStream(
	_ context.Context,
	js nats.JetStreamContext,
	subjects []string,
) error {
	// Create a new stream since it doesn't exist.
	_, err := js.AddStream(&nats.StreamConfig{
		Name:              StreamName,
		Subjects:          subjects,
		Retention:         nats.InterestPolicy, // remove acked messages
		Discard:           nats.DiscardOld,     // when the stream is full, discard old messages
		Storage:           nats.FileStorage,    // type of message storage
		MaxMsgsPerSubject: 1000,                // max stored messages per subject
		MaxMsgSize:        1 << 20,             // max single message size is 4 MB
		NoAck:             false,               // we need the "ack" system for the message queue system
	})
	// Handle potential errors during stream creation.
	if err != nil {
		return fmt.Errorf("error creating stream: %v", err)
	}
	return nil
}

// createConsumer adds a new consumer to a JetStream stream.
func createConsumer(
	js nats.JetStreamContext,
	streamName, subject, durableName string,
) error {
	_, err := js.AddConsumer(streamName, &nats.ConsumerConfig{
		Durable:        durableName,
		FilterSubject:  subject,
		AckPolicy:      nats.AckExplicitPolicy,
		DeliverSubject: subject,
	})
	if err != nil {
		return fmt.Errorf("error creating consumer: %v", err)
	}
	return nil
}
