package event

import (
	"encoding/json"
	"errors"
	"fmt"

	stan "github.com/nats-io/stan.go"
	"go.uber.org/zap"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/event"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/nats"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/index"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/model"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/service"
)

// NodeHandler interface represents the methods to handle validated and invalid nodes.
type NodeHandler interface {
	Validated() error
	ValidationFailed() error
}

// nodeHandler implements the NodeHandler interface.
type nodeHandler struct {
	svc service.NodeService
}

// NewNodeHandler returns a new NodeHandler.
func NewNodeHandler(nodeService service.NodeService) NodeHandler {
	return &nodeHandler{
		svc: nodeService,
	}
}

// Validated method listens to validated nodes from a NATS streaming server and handles them.
func (handler *nodeHandler) Validated() error {
	return event.NewNodeValidatedListener(
		nats.Client.Client(),
		index.IndexQueueGroup,
		func(msg *stan.Msg) {
			go func() {
				defer func() {
					if err := recover(); err != nil {
						logger.Error(
							fmt.Sprintf(
								"Panic occurred in nodeValidated handler: %v",
								err,
							),
							errors.New("panic"),
						)
					}
					// Acknowledge the message regardless of error.
					if err := msg.Ack(); err != nil {
						logger.Error(
							"Error when acknowledging message",
							err,
						)
					}
				}()

				var nodeValidatedData event.NodeValidatedData
				err := json.Unmarshal(msg.Data, &nodeValidatedData)
				if err != nil {
					logger.Error(
						"Error when trying to parse nodeValidatedData",
						err,
					)
					return
				}

				err = handler.svc.SetNodeValid(&model.Node{
					ProfileURL:  nodeValidatedData.ProfileURL,
					ProfileHash: &nodeValidatedData.ProfileHash,
					ProfileStr:  nodeValidatedData.ProfileStr,
					LastUpdated: &nodeValidatedData.LastUpdated,
					Version:     &nodeValidatedData.Version,
				})
				if err != nil {
					logger.Error("Failed to set node valid",
						err,
						zap.String("ProfileURL", nodeValidatedData.ProfileURL),
						zap.String("ProfileStr", nodeValidatedData.ProfileStr),
					)
					return
				}
			}()
		}).
		Listen()
}

// ValidationFailed method listens to invalid nodes from a NATS streaming server and handles them.
func (handler *nodeHandler) ValidationFailed() error {
	return event.NewNodeValidationFailedListener(
		nats.Client.Client(),
		index.IndexQueueGroup,
		func(msg *stan.Msg) {
			go func() {
				defer func() {
					if err := recover(); err != nil {
						logger.Error(
							fmt.Sprintf(
								"Panic occurred in nodeValidationFailed handler: %v",
								err,
							),
							errors.New("panic"),
						)
					}
					// Acknowledge the message regardless of error.
					if err := msg.Ack(); err != nil {
						logger.Error(
							"Error when acknowledging message",
							err,
						)
					}
				}()

				var nodeValidationFailedData event.NodeValidationFailedData
				err := json.Unmarshal(msg.Data, &nodeValidationFailedData)
				if err != nil {
					logger.Error(
						"Error when trying to parse nodeValidationFailedData",
						err,
					)
					return
				}

				err = handler.svc.SetNodeInvalid(&model.Node{
					ProfileURL:     nodeValidationFailedData.ProfileURL,
					FailureReasons: nodeValidationFailedData.FailureReasons,
					Version:        &nodeValidationFailedData.Version,
				})
				if err != nil {
					return
				}
			}()
		}).
		Listen()
}
