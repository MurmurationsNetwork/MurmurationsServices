package event

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/nats-io/nats.go"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/messaging"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/validation/internal/model"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/validation/internal/service"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/validation/internal/validation"
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
	return messaging.QueueSubscribe(
		messaging.NodeCreated,
		validation.QueueGroup,
		handler.newNodeCreatedHandler,
	)
}

// newNodeCreatedHandler handles the logic for node-created messages.
func (handler *nodeHandler) newNodeCreatedHandler(msg *nats.Msg) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error(
				fmt.Sprintf("Panic occurred in nodeCreated handler: %v", err),
				errors.New("panic"),
			)
		}
		// Acknowledge the message regardless of error.
		if err := msg.Ack(); err != nil {
			logger.Error("Error when acknowledging message", err)
		}
	}()

	var nodeCreatedData messaging.NodeCreatedData
	if err := json.Unmarshal(msg.Data, &nodeCreatedData); err != nil {
		logger.Error("Error when trying to parse nodeCreatedData", err)
		return
	}

	handler.validationService.ValidateNode(&model.Node{
		ProfileURL: nodeCreatedData.ProfileURL,
		Version:    nodeCreatedData.Version,
	})
}
