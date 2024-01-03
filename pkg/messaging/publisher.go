package messaging

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/natsclient"
)

var (
	publisherInstance *Publisher
	publisherOnce     sync.Once
)

type Publisher struct {
	natsClient *natsclient.NatsClient
}

// Publish checks for an existing Publisher instance or creates one,
// and then publishes the message to the specified subject.
func Publish(subject string, message any) error {
	var err error
	publisherOnce.Do(func() {
		publisherInstance, err = newPublisher()
	})
	if err != nil {
		return fmt.Errorf("failed to initialize publisher: %v", err)
	}

	jsonMessage, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf(
			"error marshaling message to JSON for subject '%s': %v",
			subject,
			err,
		)
	}

	return publisherInstance.publish(subject, jsonMessage)
}

func PublishSync(subject string, message any) error {
	var err error
	publisherOnce.Do(func() {
		publisherInstance, err = newPublisher()
	})
	if err != nil {
		return fmt.Errorf("failed to initialize publisher: %v", err)
	}

	jsonMessage, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf(
			"error marshaling message to JSON for subject '%s': %v",
			subject,
			err,
		)
	}

	return publisherInstance.publishSync(subject, jsonMessage)
}

// newPublisher creates a new Publisher instance.
func newPublisher() (*Publisher, error) {
	natsClient := natsclient.GetInstance()
	if natsClient == nil {
		return nil, fmt.Errorf("NATS client is not initialized")
	}
	return &Publisher{natsClient: natsClient}, nil
}

// publish publishes a message to the given subject.
func (p *Publisher) publish(subject string, message []byte) error {
	_, err := p.natsClient.JsContext.PublishAsync(subject, message)
	if err != nil {
		return fmt.Errorf(
			"failed to publish message to subject '%s': %v",
			subject,
			err,
		)
	}
	return nil
}

// publish publishes a message to the given subject.
func (p *Publisher) publishSync(subject string, message []byte) error {
	_, err := p.natsClient.JsContext.Publish(subject, message)
	if err != nil {
		return fmt.Errorf(
			"failed to publish message to subject '%s': %v",
			subject,
			err,
		)
	}
	return nil
}
