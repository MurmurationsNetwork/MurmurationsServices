package natsclient

import (
	"fmt"
	"strings"
	"sync"

	"github.com/nats-io/nats.go"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/retry"
)

// NatsClient manages NATS connections and subscriptions.
type NatsClient struct {
	conn          *nats.Conn
	JsContext     nats.JetStreamContext
	subscriptions []*nats.Subscription
}

var (
	instance *NatsClient
	once     sync.Once
)

// Initialize sets up the NATS client with the provided URL.
func Initialize(url string) error {
	var err error
	once.Do(func() {
		instance = &NatsClient{}
		instance.conn, err = connectToNATS(url)
		if err == nil {
			instance.JsContext, err = instance.conn.JetStream()
		}
		if err == nil {
			err = instance.ensureStreamExists()
		}
	})
	return err
}

// GetInstance retrieves the initialized NatsClient instance.
func GetInstance() *NatsClient {
	if instance == nil || instance.conn == nil {
		panic("NATS client is not initialized. Call Initialize first.")
	}
	return instance
}

// AddSubscription adds a subscription to the NatsClient for management.
func (c *NatsClient) AddSubscription(sub *nats.Subscription) {
	c.subscriptions = append(c.subscriptions, sub)
}

// Disconnect gracefully closes the NATS connection and drains subscriptions.
func (c *NatsClient) Disconnect() error {
	var errStrings []string

	if err := c.drainSubscriptions(); err != nil {
		errStrings = append(
			errStrings,
			fmt.Sprintf("error draining subscriptions: %v", err),
		)
	}

	if c.conn != nil {
		if err := c.conn.Drain(); err != nil {
			errStrings = append(
				errStrings,
				fmt.Sprintf("error draining connection: %v", err),
			)
		}
		c.conn.Close()
	}

	if len(errStrings) > 0 {
		return fmt.Errorf(
			"disconnect encountered issues: %s",
			strings.Join(errStrings, ", "),
		)
	}

	return nil
}

// connectToNATS establishes a connection to the NATS server at the given URL.
func connectToNATS(natsURL string) (*nats.Conn, error) {
	var conn *nats.Conn
	err := retry.Do(func() error {
		var err error
		conn, err = nats.Connect(natsURL)
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}
	return conn, nil
}

// ensureStreamExists ensures the required stream exists in NATS.
func (c *NatsClient) ensureStreamExists() error {
	_, err := c.JsContext.StreamInfo(streamName)
	if err == nil {
		return nil
	}
	if err != nats.ErrStreamNotFound {
		return fmt.Errorf("error checking stream existence: %v", err)
	}
	return c.createStream()
}

// createStream creates a new stream in NATS JetStream.
func (c *NatsClient) createStream() error {
	streamConfig := &nats.StreamConfig{
		Name:              streamName,
		Subjects:          []string{"NODES.>"},
		Retention:         nats.WorkQueuePolicy,
		Discard:           nats.DiscardOld,
		Storage:           nats.FileStorage,
		MaxMsgsPerSubject: 1000,
		MaxMsgSize:        1 << 20, // 1 MB
		NoAck:             false,
	}

	info, err := c.JsContext.AddStream(streamConfig)
	if err != nil {
		return fmt.Errorf("error creating stream: %v", err)
	}

	logger.Info(fmt.Sprintf("Stream created: %+v\n", info))

	return nil
}

// drainSubscriptions drains all managed subscriptions.
func (c *NatsClient) drainSubscriptions() error {
	var errStrings []string
	for _, sub := range c.subscriptions {
		if err := sub.Drain(); err != nil {
			errStrings = append(
				errStrings,
				fmt.Sprintf("error draining subscription: %v", err),
			)
		}
	}
	if len(errStrings) > 0 {
		return fmt.Errorf(
			"issues draining subscriptions: %s",
			strings.Join(errStrings, ", "),
		)
	}
	return nil
}
