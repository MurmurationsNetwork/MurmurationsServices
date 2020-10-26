package events

import (
	stan "github.com/nats-io/stan.go"
)

func NewNodeCreatedPublisher(client stan.Conn) Publisher {
	return NewPublisher(&publisherConfig{
		Client:  client,
		Subject: nodeCreated,
	})
}

func NewNodeValidatedPublisher(client stan.Conn) Publisher {
	return NewPublisher(&publisherConfig{
		Client:  client,
		Subject: nodeValidated,
	})
}

func NewNodeValidationFailedPublisher(client stan.Conn) Publisher {
	return NewPublisher(&publisherConfig{
		Client:  client,
		Subject: nodeValidationFailed,
	})
}
