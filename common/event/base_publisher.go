package event

import (
	"encoding/json"
	"os"

	"github.com/MurmurationsNetwork/MurmurationsServices/common/logger"
	stan "github.com/nats-io/stan.go"
)

type Publisher interface {
	Publish(data interface{})
	SetAckHandler(ackHandler stan.AckHandler)
}

type publisherConfig struct {
	Client  stan.Conn
	Subject Subject
}

type publisher struct {
	client     stan.Conn
	subject    Subject
	ackHandler stan.AckHandler
}

func NewPublisher(config *publisherConfig) Publisher {
	return &publisher{
		client:  config.Client,
		subject: config.Subject,
		ackHandler: func(guid string, err error) {
			if err != nil {
				logger.Error("error when trying to publish "+string(config.Subject)+" event.", err)
			}
		},
	}
}

func (p *publisher) SetAckHandler(ackHandler stan.AckHandler) {
	p.ackHandler = ackHandler
}

func (p *publisher) Publish(data interface{}) {
	// FIXME: Use Abstraction
	if os.Getenv("ENV") == "test" {
		return
	}
	msg, _ := json.Marshal(data)
	p.client.PublishAsync(string(p.subject), msg, p.ackHandler)
}
