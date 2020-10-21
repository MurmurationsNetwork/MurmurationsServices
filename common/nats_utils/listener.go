package nats_utils

import (
	"fmt"
	"time"

	stan "github.com/nats-io/stan.go"
)

type Listener interface {
	Listen() error
}

type listener struct {
	client      stan.Conn
	subject     string
	qgroup      string
	maxInflight int
	ackWait     time.Duration
	onMessage   stan.MsgHandler
}

func NewListener(client stan.Conn, subject string, qgroup string) Listener {
	return &listener{
		client:      client,
		subject:     subject,
		qgroup:      qgroup,
		maxInflight: 50,
		ackWait:     10 * time.Second,
		onMessage: func(msg *stan.Msg) {
			fmt.Println("receiving message", msg.Sequence, string(msg.Data))
			msg.Ack()
		},
	}
}

func (l *listener) subscriptionOptions() []stan.SubscriptionOption {
	return []stan.SubscriptionOption{
		stan.SetManualAckMode(),
		stan.DeliverAllAvailable(),
		stan.DurableName(l.qgroup),
		stan.MaxInflight(l.maxInflight),
		stan.AckWait(l.ackWait),
	}
}

func (l *listener) Listen() error {
	_, err := l.client.QueueSubscribe(l.subject, l.qgroup, l.onMessage, l.subscriptionOptions()...)
	if err != nil {
		return err
	}
	return nil
}
