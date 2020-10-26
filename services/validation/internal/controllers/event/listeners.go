package event

import (
	"encoding/json"

	"github.com/MurmurationsNetwork/MurmurationsServices/common/events"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/validation/internal/datasources/nats"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/validation/internal/services"
	"github.com/nats-io/stan.go"
)

var HandleNodeCreated = events.NewNodeCreatedListener(nats.Client(), qgroup, func(msg *stan.Msg) {
	var nodeCreatedData events.NodeCreatedData
	err := json.Unmarshal(msg.Data, &nodeCreatedData)
	if err != nil {
		logger.Error("error when trying to parsing nodeCreatedData", err)
		return
	}

	err = services.ValidationService.ValidateNode(nodeCreatedData.ProfileUrl, nodeCreatedData.LinkedSchemas)
	if err != nil {
		return
	}

	msg.Ack()
})
