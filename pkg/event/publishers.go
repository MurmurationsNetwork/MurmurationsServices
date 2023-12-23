package event

import "github.com/nats-io/nats.go"

func NewNodeCreatedPublisher(js nats.JetStreamContext) Publisher {
	return NewPublisher(&publisherConfig{
		JetStream: js,
		Subject:   NodeCreated,
	})
}

func NewNodeValidatedPublisher(js nats.JetStreamContext) Publisher {
	return NewPublisher(&publisherConfig{
		JetStream: js,
		Subject:   NodeValidated,
	})
}

func NewNodeValidationFailedPublisher(js nats.JetStreamContext) Publisher {
	return NewPublisher(&publisherConfig{
		JetStream: js,
		Subject:   NodeValidationFailed,
	})
}
