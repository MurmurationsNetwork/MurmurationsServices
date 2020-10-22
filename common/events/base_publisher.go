package events

import (
	"encoding/json"
	"fmt"

	stan "github.com/nats-io/stan.go"
)

type Publisher interface {
	Publish(data interface{}) error
}

type PublisherConfig struct {
	Client     stan.Conn
	Subject    Subject
	AckHandler stan.AckHandler
}

type publisher struct {
	client     stan.Conn
	subject    Subject
	AckHandler stan.AckHandler
}

func NewPublisher(config *PublisherConfig) Publisher {
	return &publisher{
		client:     config.Client,
		subject:    config.Subject,
		AckHandler: config.AckHandler,
	}
}

func DefaultAckHandler() stan.MsgHandler {
	return func(msg *stan.Msg) {
		fmt.Println("receiving message: ", msg.Sequence, string(msg.Data))
		msg.Ack()
	}
}

func (p *publisher) Publish(data interface{}) error {
	msg, _ := json.Marshal(data)
	_, err := p.client.PublishAsync(string(p.subject), msg, p.AckHandler)
	if err != nil {
		return err
	}
	return nil
}
