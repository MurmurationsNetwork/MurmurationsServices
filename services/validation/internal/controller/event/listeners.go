package event

import (
	"encoding/json"
	"fmt"

	"github.com/MurmurationsNetwork/MurmurationsServices/common/event"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/validation/internal/datasource/nats"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/validation/internal/domain/node"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/validation/internal/service"
	"github.com/nats-io/stan.go"
)

var HandleNodeCreated = event.NewNodeCreatedListener(nats.Client(), qgroup, func(msg *stan.Msg) {
	var nodeCreatedData event.NodeCreatedData
	err := json.Unmarshal(msg.Data, &nodeCreatedData)
	if err != nil {
		logger.Error("error when trying to parsing nodeCreatedData", err)
		return
	}

	fmt.Println("==================================")
	fmt.Printf("nodeCreatedData %+v \n", nodeCreatedData)
	fmt.Println("==================================")

	service.ValidationService.ValidateNode(&node.Node{
		ProfileURL: nodeCreatedData.ProfileURL,
		Version:    nodeCreatedData.Version,
	})

	msg.Ack()
})
