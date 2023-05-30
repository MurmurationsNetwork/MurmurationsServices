package event

import (
	stan "github.com/nats-io/stan.go"
)

func NewNodeCreatedListener(
	client stan.Conn,
	qgroup string,
	handler stan.MsgHandler,
) Listener {
	return NewListener(&ListenerConfig{
		Client:     client,
		Subject:    nodeCreated,
		Qgroup:     qgroup,
		MsgHandler: handler,
	})
}

func NewNodeValidatedListener(
	client stan.Conn,
	qgroup string,
	handler stan.MsgHandler,
) Listener {
	return NewListener(&ListenerConfig{
		Client:     client,
		Subject:    nodeValidated,
		Qgroup:     qgroup,
		MsgHandler: handler,
	})
}

func NewNodeValidationFailedListener(
	client stan.Conn,
	qgroup string,
	handler stan.MsgHandler,
) Listener {
	return NewListener(&ListenerConfig{
		Client:     client,
		Subject:    nodeValidationFailed,
		Qgroup:     qgroup,
		MsgHandler: handler,
	})
}
