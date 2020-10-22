package publishers

import (
	"fmt"

	"github.com/MurmurationsNetwork/MurmurationsServices/common/events"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/datasources/nats"
)

var NodeCreated = events.NewPublisher(&events.PublisherConfig{
	Client:  nats.Client(),
	Subject: events.NodeCreated,
	AckHandler: func(guid string, err error) {
		if err != nil {
			fmt.Printf("Event publish failed.")
			return
		}
		fmt.Printf("Event published!")
	},
})
