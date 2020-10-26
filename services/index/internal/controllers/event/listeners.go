package event

import (
	"encoding/json"
	"fmt"

	"github.com/MurmurationsNetwork/MurmurationsServices/common/events"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/datasources/nats"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/domain/nodes"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/services"
	"github.com/nats-io/stan.go"
)

var HandleNodeValidated = events.NewNodeValidatedListener(nats.Client(), qgroup, func(msg *stan.Msg) {
	var nodeValidatedData events.NodeValidatedData
	err := json.Unmarshal(msg.Data, &nodeValidatedData)
	if err != nil {
		logger.Error("error when trying to parse nodeValidatedData", err)
		return
	}

	fmt.Println("==================================")
	fmt.Printf("nodeValidatedData %+v \n", nodeValidatedData)
	fmt.Println("==================================")

	err = services.NodeService.SetNodeValid(nodes.Node{
		ProfileUrl:    nodeValidatedData.ProfileUrl,
		LastValidated: nodeValidatedData.LastValidated,
	})
	if err != nil {
		return
	}

	msg.Ack()
})

var HandleNodeValidationFailed = events.NewNodeValidationFailedListener(nats.Client(), qgroup, func(msg *stan.Msg) {
	var nodeValidationFailedData events.NodeValidationFailedData
	err := json.Unmarshal(msg.Data, &nodeValidationFailedData)
	if err != nil {
		logger.Error("error when trying to parse nodeValidationFailedData", err)
		return
	}

	fmt.Println("==================================")
	fmt.Printf("nodeValidationFailedData %+v \n", nodeValidationFailedData)
	fmt.Println("==================================")

	err = services.NodeService.SetNodeInValid(nodes.Node{
		ProfileUrl:    nodeValidationFailedData.ProfileUrl,
		FailedReasons: &nodeValidationFailedData.FailedReasons,
	})
	if err != nil {
		return
	}

	msg.Ack()
})
