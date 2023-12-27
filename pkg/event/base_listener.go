package event

import (
	"errors"
	"time"

	"github.com/nats-io/nats.go"
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
	UpdateOptions(opts ...nats.SubOpt)
}

type ListenerConfig struct {
	JetStream  nats.JetStreamContext
	Subject    Subject
	Qgroup     string
	MsgHandler nats.MsgHandler
}

type listener struct {
	js         nats.JetStreamContext
	subject    Subject
	qgroup     string
	msgHandler nats.MsgHandler
	opts       []nats.SubOpt
}

func NewListener(config *ListenerConfig) Listener {
	return &listener{
		js:      config.JetStream,
		subject: config.Subject,
		// A queue group (qgroup) in NATS is a mechanism for load balancing
		// messages among multiple subscribers by ensuring that each message
		// is delivered to only one subscriber in the group,
		qgroup:     config.Qgroup,
		msgHandler: config.MsgHandler,
		opts: []nats.SubOpt{
			nats.Durable(string(config.Subject) + "_consumer"),
			nats.MaxAckPending(DefaultMaxInflight),
			nats.AckWait(DefaultAckWait),
		},
	}
}

// UpdateOptions overrides the default options.
func (l *listener) UpdateOptions(opts ...nats.SubOpt) {
	l.opts = append(l.opts, opts...)
}

func (l *listener) Listen() error {
	if l.msgHandler == nil {
		return ErrNilMsgHandler
	}

	_, err := l.js.QueueSubscribe(
		string(l.subject),
		l.qgroup,
		l.msgHandler,
		l.opts...,
	)

	return err
}
