package messaging

import (
	"fmt"
	"sync"

	"github.com/nats-io/nats.go"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/natsclient"
)

type MessageHandler func(msg *nats.Msg)

var (
	subscriberInstance *Subscriber
	subscriberOnce     sync.Once
)

type Subscriber struct {
	natsClient *natsclient.NatsClient
}

// QueueSubscribe checks for an existing Subscriber instance or creates one,
// and then subscribes to the specified queue.
func QueueSubscribe(
	subject, queue string,
	handler MessageHandler,
) error {
	var err error
	subscriberOnce.Do(func() {
		subscriberInstance, err = newSubscriber()
	})
	if err != nil {
		return err
	}

	return subscriberInstance.queueSubscribe(subject, queue, handler)
}

// newSubscriber creates a new Subscriber instance.
func newSubscriber() (*Subscriber, error) {
	natsClient := natsclient.GetInstance()
	if natsClient == nil {
		return nil, fmt.Errorf("NATS client is not initialized")
	}
	return &Subscriber{natsClient: natsClient}, nil
}

// queueSubscribe subscribes to the given subject and queue.
func (s *Subscriber) queueSubscribe(
	subject, queue string,
	handler MessageHandler,
) error {
	sub, err := s.natsClient.JsContext.QueueSubscribe(
		subject,
		queue,
		func(msg *nats.Msg) {
			handler(msg)
		},
		nats.Durable(subject+"_consumer"),
	)
	if err != nil {
		return err
	}

	s.natsClient.AddSubscription(sub)
	return nil
}
