package nats_utils

import (
	"fmt"
	"time"

	stan "github.com/nats-io/stan.go"
)

const (
	DefaultAckWait     = 20 * time.Second
	DefaultMaxInflight = 50
)

type Listener interface {
	Listen() error
	UpdateOptions(opts ...stan.SubscriptionOption)
}

type listener struct {
	client    stan.Conn
	subject   string
	qgroup    string
	onMessage stan.MsgHandler
	opts      []stan.SubscriptionOption
}

func NewListener(client stan.Conn, subject string, qgroup string) Listener {
	return &listener{
		client:    client,
		subject:   subject,
		qgroup:    qgroup,
		onMessage: defaultOnMessage(),
		opts:      defaultSubscriptionOptions(qgroup),
	}
}

func defaultOnMessage() stan.MsgHandler {
	return func(msg *stan.Msg) {
		fmt.Println("receiving message", msg.Sequence, string(msg.Data))
		msg.Ack()
	}
}

func defaultSubscriptionOptions(qgroup string) []stan.SubscriptionOption {
	return []stan.SubscriptionOption{
		stan.SetManualAckMode(),
		stan.DeliverAllAvailable(),
		stan.DurableName(qgroup),
		stan.MaxInflight(DefaultMaxInflight),
		stan.AckWait(DefaultAckWait),
	}
}

func (l *listener) Listen() error {
	_, err := l.client.QueueSubscribe(l.subject, l.qgroup, l.onMessage, l.opts...)
	if err != nil {
		return err
	}
	return nil
}

func (l *listener) UpdateOptions(opts ...stan.SubscriptionOption) {
	l.opts = append(l.opts, opts...)
}
