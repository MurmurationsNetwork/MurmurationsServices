package event

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/MurmurationsNetwork/MurmurationsServices/common/event"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/nats"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/domain/node"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/service"
	"github.com/nats-io/stan.go"
)

type NodeHandler interface {
	Validated() error
	ValidationFailed() error
}

type nodeHandler struct {
	nodeService service.NodesService
}

func NewNodeHandler(nodeService service.NodesService) NodeHandler {
	return &nodeHandler{
		nodeService: nodeService,
	}
}

func (handler *nodeHandler) Validated() error {
	return event.NewNodeValidatedListener(nats.Client.Client(), QGROOP, func(msg *stan.Msg) {
		go func() {
			defer func() {
				if err := recover(); err != nil {
					logger.Error(fmt.Sprintf("Panic occurred in nodeValidated handler: %v", err), errors.New("panic"))
				}
			}()

			var nodeValidatedData event.NodeValidatedData
			err := json.Unmarshal(msg.Data, &nodeValidatedData)
			if err != nil {
				logger.Error("error when trying to parse nodeValidatedData", err)
				return
			}

			err = handler.nodeService.SetNodeValid(&node.Node{
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
		}()
	}).Listen()
}

func (handler *nodeHandler) ValidationFailed() error {
	return event.NewNodeValidationFailedListener(nats.Client.Client(), QGROOP, func(msg *stan.Msg) {
		go func() {
			defer func() {
				if err := recover(); err != nil {
					logger.Error(fmt.Sprintf("Panic occurred in nodeValidationFailed handler: %v", err), errors.New("panic"))
				}
			}()

			var nodeValidationFailedData event.NodeValidationFailedData
			err := json.Unmarshal(msg.Data, &nodeValidationFailedData)
			if err != nil {
				logger.Error("error when trying to parse nodeValidationFailedData", err)
				return
			}

			err = handler.nodeService.SetNodeInvalid(&node.Node{
				ProfileURL:     nodeValidationFailedData.ProfileURL,
				FailureReasons: &nodeValidationFailedData.FailureReasons,
				Version:        &nodeValidationFailedData.Version,
			})
			if err != nil {
				return
			}

			msg.Ack()
		}()
	}).Listen()
}
