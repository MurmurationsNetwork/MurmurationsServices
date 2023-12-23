package nats

import "github.com/nats-io/nats.go"

type natsClientInterface interface {
	JetStream() nats.JetStreamContext
	Disconnect()
	setClient(*nats.Conn, nats.JetStreamContext)
}
