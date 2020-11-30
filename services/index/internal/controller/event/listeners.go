package event

import (
	"encoding/json"

	"github.com/MurmurationsNetwork/MurmurationsServices/common/event"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/nats"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/domain/node"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/service"
	"github.com/nats-io/stan.go"
)

func HandleNodeValidated() event.Listener {
	return event.NewNodeValidatedListener(nats.Client.Client(), qgroup, func(msg *stan.Msg) {
		var nodeValidatedData event.NodeValidatedData
		err := json.Unmarshal(msg.Data, &nodeValidatedData)
		if err != nil {
			logger.Error("error when trying to parse nodeValidatedData", err)
			return
		}

		err = service.NodeService.SetNodeValid(node.Node{
			ProfileURL:    nodeValidatedData.ProfileURL,
			ProfileHash:   &nodeValidatedData.ProfileHash,
			ProfileStr:    nodeValidatedData.ProfileStr,
			LastValidated: &nodeValidatedData.LastValidated,
			Version:       &nodeValidatedData.Version,
		})
		if err != nil {
			return
		}

		msg.Ack()
	})
}

func HandleNodeValidationFailed() event.Listener {
	return event.NewNodeValidationFailedListener(nats.Client.Client(), qgroup, func(msg *stan.Msg) {
		var nodeValidationFailedData event.NodeValidationFailedData
		err := json.Unmarshal(msg.Data, &nodeValidationFailedData)
		if err != nil {
			logger.Error("error when trying to parse nodeValidationFailedData", err)
			return
		}

		err = service.NodeService.SetNodeInvalid(&node.Node{
			ProfileURL:     nodeValidationFailedData.ProfileURL,
			FailureReasons: &nodeValidationFailedData.FailureReasons,
			Version:        &nodeValidationFailedData.Version,
		})
		if err != nil {
			return
		}

		msg.Ack()
	})

}
