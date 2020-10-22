package listeners

import (
	"encoding/json"
	"fmt"

	"github.com/MurmurationsNetwork/MurmurationsServices/common/events"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/validation/internal/datasources/nats"
	"github.com/nats-io/stan.go"
)

const QGroup = "validation-svc-qgroup"

var NodeCreated = events.NewListener(&events.ListenerConfig{
	Client:  nats.Client(),
	Subject: events.NodeCreated,
	QGroup:  QGroup,
	OnMessage: func(msg *stan.Msg) {
		var data events.NodeCreatedData
		json.Unmarshal(msg.Data, &data)
		fmt.Printf("%+v %+v \n", msg.Sequence, data)
		msg.Ack()
	},
})
