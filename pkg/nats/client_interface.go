package nats

import (
	"github.com/nats-io/nats.go"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/event"
)

// natsClientInterface defines the interface for a NATS client.
type natsClientInterface interface {
	// JetStream returns the JetStream context of the NATS client.
	JetStream() nats.JetStreamContext

	// SubscribeToSubjects allows the client to subscribe to a
	// variable number of subjects.
	SubscribeToSubjects(subjects ...event.Subject) error

	// Disconnect gracefully closes the connection to the NATS server.
	Disconnect()
}
