package event

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/nats-io/nats.go"
)

// Publisher defines the interface for publishing messages.
type Publisher interface {
	Publish(data interface{}) error
	PublishSync(data interface{}) error
}

// publisherConfig holds the configuration for creating a new publisher.
type publisherConfig struct {
	JetStream nats.JetStreamContext // Context for NATS JetStream.
	Subject   Subject               // Subject under which messages will be published.
	Stream    string                // Name of the stream in NATS JetStream.
}

// publisher implements the Publisher interface using NATS JetStream.
type publisher struct {
	js      nats.JetStreamContext // NATS JetStream context.
	subject Subject               // Subject for publishing messages.
}

// NewPublisher creates a new publisher with the given configuration.
func NewPublisher(config *publisherConfig) Publisher {
	return &publisher{
		js:      config.JetStream,
		subject: config.Subject,
	}
}

// Publish sends a message asynchronously.
func (p *publisher) Publish(data interface{}) error {
	if os.Getenv("APP_ENV") == "test" {
		// In a test environment, skip actual publishing.
		return nil
	}

	msg, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to serialize data into JSON: %w", err)
	}

	_, err = p.js.PublishAsync(string(p.subject), msg)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	return nil
}

// PublishSync sends a message synchronously.
func (p *publisher) PublishSync(data interface{}) error {
	if os.Getenv("APP_ENV") == "test" {
		// In a test environment, skip actual publishing.
		return nil
	}

	msg, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to serialize data into JSON: %w", err)
	}

	_, err = p.js.Publish(string(p.subject), msg)
	if err != nil {
		return fmt.Errorf("failed to publish message synchronously: %w", err)
	}

	return nil
}
