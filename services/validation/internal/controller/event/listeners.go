package event

import (
	"encoding/json"
	"errors"
	"fmt"

	stan "github.com/nats-io/stan.go"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/event"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/nats"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/validation/internal/model"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/validation/internal/service"
)

// NodeHandler provides an interface for handling node events.
type NodeHandler interface {
	NewNodeCreatedListener() error
}

type nodeHandler struct {
	validationService service.ValidationService
}

// NewNodeHandler creates a new NodeHandler with the provided validation service.
func NewNodeHandler(validationService service.ValidationService) NodeHandler {
	return &nodeHandler{
		validationService: validationService,
	}
}

// NewNodeCreatedListener starts a listener for node-created events.
func (handler *nodeHandler) NewNodeCreatedListener() error {
	return event.NewNodeCreatedListener(
		nats.Client.Client(),
		qgroup,
		func(msg *stan.Msg) {
			go func() {
				defer func() {
					if err := recover(); err != nil {
						logger.Error(
							fmt.Sprintf(
								"Panic occurred in nodeCreated handler: %v",
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

				var nodeCreatedData event.NodeCreatedData
				err := json.Unmarshal(msg.Data, &nodeCreatedData)
				if err != nil {
					logger.Error(
						"Error when trying to parse nodeCreatedData",
						err,
					)
					return
				}

				handler.validationService.ValidateNode(&model.Node{
					ProfileURL: nodeCreatedData.ProfileURL,
					Version:    nodeCreatedData.Version,
				})
			}()
		},
	).Listen()
}
