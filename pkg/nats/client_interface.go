package nats

import (
	"github.com/nats-io/nats.go"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/event"
)

type natsClientInterface interface {
	JetStream() nats.JetStreamContext
	SubscribeToSubjects(subjects ...event.Subject) error
	Disconnect()
	setClient(*nats.Conn, nats.JetStreamContext)
}
