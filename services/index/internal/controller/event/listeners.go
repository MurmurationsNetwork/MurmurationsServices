package event

import (
	"encoding/json"
	"fmt"

	"github.com/MurmurationsNetwork/MurmurationsServices/common/event"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/datasources/nats"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/domain/node"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/service"
	"github.com/nats-io/stan.go"
)

var HandleNodeValidated = event.NewNodeValidatedListener(nats.Client(), qgroup, func(msg *stan.Msg) {
	var nodeValidatedData event.NodeValidatedData
	err := json.Unmarshal(msg.Data, &nodeValidatedData)
	if err != nil {
		logger.Error("error when trying to parse nodeValidatedData", err)
		return
	}

	fmt.Println("==================================")
	fmt.Printf("nodeValidatedData %+v \n", nodeValidatedData)
	fmt.Println("==================================")

	err = service.NodeService.SetNodeValid(node.Node{
		ProfileUrl:    nodeValidatedData.ProfileUrl,
		ProfileHash:   &nodeValidatedData.ProfileHash,
		LastChecked: nodeValidatedData.LastChecked,
		Version:       &nodeValidatedData.Version,
	})
	if err != nil {
		return
	}

	msg.Ack()
})

var HandleNodeValidationFailed = event.NewNodeValidationFailedListener(nats.Client(), qgroup, func(msg *stan.Msg) {
	var nodeValidationFailedData event.NodeValidationFailedData
	err := json.Unmarshal(msg.Data, &nodeValidationFailedData)
	if err != nil {
		logger.Error("error when trying to parse nodeValidationFailedData", err)
		return
	}

	fmt.Println("==================================")
	fmt.Printf("nodeValidationFailedData %+v \n", nodeValidationFailedData)
	fmt.Println("==================================")

	err = service.NodeService.SetNodeInvalid(node.Node{
		ProfileUrl:    nodeValidationFailedData.ProfileUrl,
		FailedReasons: &nodeValidationFailedData.FailedReasons,
		Version:       &nodeValidationFailedData.Version,
	})
	if err != nil {
		return
	}

	msg.Ack()
})
