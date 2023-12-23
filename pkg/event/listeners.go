package event

import (
	"github.com/nats-io/nats.go"
)

func NewNodeCreatedListener(
	js nats.JetStreamContext,
	qgroup string,
	handler nats.MsgHandler,
) Listener {
	return NewListener(&ListenerConfig{
		JetStream:  js,
		Subject:    NodeCreated,
		Qgroup:     qgroup,
		MsgHandler: handler,
	})
}

func NewNodeValidatedListener(
	js nats.JetStreamContext,
	qgroup string,
	handler nats.MsgHandler,
) Listener {
	return NewListener(&ListenerConfig{
		JetStream:  js,
		Subject:    NodeValidated,
		Qgroup:     qgroup,
		MsgHandler: handler,
	})
}

func NewNodeValidationFailedListener(
	js nats.JetStreamContext,
	qgroup string,
	handler nats.MsgHandler,
) Listener {
	return NewListener(&ListenerConfig{
		JetStream:  js,
		Subject:    NodeValidationFailed,
		Qgroup:     qgroup,
		MsgHandler: handler,
	})
}
