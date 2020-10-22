package events

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

type ListenerConfig struct {
	Client    stan.Conn
	Subject   Subject
	QGroup    string
	OnMessage stan.MsgHandler
}

type listener struct {
	client    stan.Conn
	subject   Subject
	qgroup    string
	onMessage stan.MsgHandler
	opts      []stan.SubscriptionOption
}

func NewListener(config *ListenerConfig) Listener {
	return &listener{
		client:    config.Client,
		subject:   config.Subject,
		qgroup:    config.QGroup,
		onMessage: config.OnMessage,
		opts:      defaultSubscriptionOptions(config.QGroup),
	}
}

func DefaultOnMessage() stan.MsgHandler {
	return func(msg *stan.Msg) {
		fmt.Println("receiving message: ", msg.Sequence, string(msg.Data))
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

// UpdateOptions overrides the default options.
func (l *listener) UpdateOptions(opts ...stan.SubscriptionOption) {
	l.opts = append(l.opts, opts...)
}

func (l *listener) Listen() error {
	_, err := l.client.QueueSubscribe(string(l.subject), l.qgroup, l.onMessage, l.opts...)
	if err != nil {
		return err
	}
	return nil
}
