package event

import (
	"encoding/json"
	"fmt"
	"os"

	stan "github.com/nats-io/stan.go"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/logger"
)

type Publisher interface {
	Publish(data interface{})
	PublishSync(data interface{}) error
	SetAckHandler(ackHandler stan.AckHandler)
}

type publisherConfig struct {
	Client  stan.Conn
	Subject Subject
}

type publisher struct {
	client     stan.Conn
	subject    Subject
	ackHandler stan.AckHandler
}

func NewPublisher(config *publisherConfig) Publisher {
	return &publisher{
		client:  config.Client,
		subject: config.Subject,
		ackHandler: func(guid string, err error) {
			if err != nil {
				logger.Error(
					"error when trying to publish "+string(
						config.Subject,
					)+" event.",
					err,
				)
			}
		},
	}
}

func (p *publisher) SetAckHandler(ackHandler stan.AckHandler) {
	p.ackHandler = ackHandler
}

func (p *publisher) Publish(data interface{}) {
	// FIXME: Use Abstraction
	if os.Getenv("APP_ENV") == "test" {
		return
	}
	msg, _ := json.Marshal(data)
	_, _ = p.client.PublishAsync(string(p.subject), msg, p.ackHandler)
}

// PublishSync sends a message to the designated subject using a synchronous approach.
func (p *publisher) PublishSync(data interface{}) error {
	if os.Getenv("APP_ENV") == "test" {
		return nil
	}

	msg, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to serialize data into JSON: %w", err)
	}

	err = p.client.Publish(string(p.subject), msg)
	if err != nil {
		return fmt.Errorf(
			"failed to publish message to subject %s: %w",
			p.subject,
			err,
		)
	}

	return nil
}
