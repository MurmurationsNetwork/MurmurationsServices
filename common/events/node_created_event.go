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

func NewNodeCreatedListener(client stan.Conn, qgroup string, handler stan.MsgHandler) Listener {
	return NewListener(&ListenerConfig{
		Client:     client,
		Subject:    nodeCreated,
		Qgroup:     qgroup,
		MsgHandler: handler,
	})
}
