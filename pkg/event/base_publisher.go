package event

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/nats-io/nats.go"
)

type Publisher interface {
	Publish(data interface{}) error
	PublishSync(data interface{}) error
}

type publisherConfig struct {
	JetStream nats.JetStreamContext
	Subject   Subject
	Stream    string
}

type publisher struct {
	js      nats.JetStreamContext
	subject Subject
}

func NewPublisher(config *publisherConfig) Publisher {
	return &publisher{
		js:      config.JetStream,
		subject: config.Subject,
	}
}

func (p *publisher) Publish(data interface{}) error {
	if os.Getenv("APP_ENV") == "test" {
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

func (p *publisher) PublishSync(data interface{}) error {
	if os.Getenv("APP_ENV") == "test" {
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
