package event

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/MurmurationsNetwork/MurmurationsServices/common/event"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/nats"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/validation/internal/domain/node"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/validation/internal/service"
	"github.com/nats-io/stan.go"
)

type NodeHandler interface {
	NewNodeCreatedListener() error
}

type nodeHandler struct {
	validationService service.ValidationService
}

func NewNodeHandler(validationService service.ValidationService) NodeHandler {
	return &nodeHandler{
		validationService: validationService,
	}
}

func (handler *nodeHandler) NewNodeCreatedListener() error {
	return event.NewNodeCreatedListener(nats.Client.Client(), qgroup, func(msg *stan.Msg) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error(fmt.Sprintf("Panic occurred in nodeCreated handler: %v", err), errors.New("panic"))
			}
		}()

		var nodeCreatedData event.NodeCreatedData
		err := json.Unmarshal(msg.Data, &nodeCreatedData)
		if err != nil {
			logger.Error("Error when trying to parsing nodeCreatedData", err)
			return
		}

		handler.validationService.ValidateNode(&node.Node{
			ProfileURL: nodeCreatedData.ProfileURL,
			Version:    nodeCreatedData.Version,
		})

		msg.Ack()
	}).Listen()
}
