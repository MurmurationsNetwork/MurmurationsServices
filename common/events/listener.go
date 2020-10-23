package events

import (
	"errors"
	"fmt"
	"time"

	stan "github.com/nats-io/stan.go"
)

const (
	DefaultAckWait     = 20 * time.Second
	DefaultMaxInflight = 50
)

var (
	ErrNilMsgHandler = errors.New("listener: message handler cannot be nil")
)

type Listener interface {
	Listen() error
	UpdateOptions(opts ...stan.SubscriptionOption)
}

type ListenerConfig struct {
	Client     stan.Conn
	Subject    Subject
	Qgroup     string
	MsgHandler stan.MsgHandler
}

type listener struct {
	client     stan.Conn
	subject    Subject
	qgroup     string
	msgHandler stan.MsgHandler
	opts       []stan.SubscriptionOption
}

func NewListener(config *ListenerConfig) Listener {
	return &listener{
		client:     config.Client,
		subject:    config.Subject,
		qgroup:     config.Qgroup,
		msgHandler: config.MsgHandler,
		opts: []stan.SubscriptionOption{
			stan.SetManualAckMode(),
			stan.DeliverAllAvailable(),
			stan.DurableName(config.Qgroup),
			stan.MaxInflight(DefaultMaxInflight),
			stan.AckWait(DefaultAckWait),
		},
	}
}

// UpdateOptions overrides the default options.
func (l *listener) UpdateOptions(opts ...stan.SubscriptionOption) {
	l.opts = append(l.opts, opts...)
}

func DefaultMsgHandler() stan.MsgHandler {
	return func(msg *stan.Msg) {
		fmt.Println("receiving message: ", msg.Sequence, string(msg.Data))
		msg.Ack()
	}
}

func (l *listener) Listen() error {
	if l.msgHandler == nil {
		return ErrNilMsgHandler
	}

	_, err := l.client.QueueSubscribe(string(l.subject), l.qgroup, l.msgHandler, l.opts...)
	if err != nil {
		return err
	}
	return nil
}
