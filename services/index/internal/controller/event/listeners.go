package event

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync/atomic"

	natsio "github.com/nats-io/nats.go"
	"go.uber.org/zap"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/messaging"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/index"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/model"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/service"
)

// NodeHandler defines methods for handling validated and invalid node events.
type NodeHandler interface {
	Validated() error
	ValidationFailed() error
}

// nodeHandler handles node-related events.
type nodeHandler struct {
	svc service.NodeService
}

// NewNodeHandler creates a new handler for node-related events.
func NewNodeHandler(nodeService service.NodeService) NodeHandler {
	return &nodeHandler{svc: nodeService}
}

// Validated sets up a listener for validated node events and processes them.
func (handler *nodeHandler) Validated() error {
	err := messaging.QueueSubscribe(
		messaging.NodeValidated,
		index.QueueGroup,
		handler.processValidatedNode,
	)
	if err != nil {
		return fmt.Errorf(
			"failed to subscribe to '%s': %v",
			messaging.NodeValidated,
			err,
		)
	}

	return nil
}

// ValidationFailed sets up a listener for validation failed node events and
// processes them.
func (handler *nodeHandler) ValidationFailed() error {
	err := messaging.QueueSubscribe(
		messaging.NodeValidationFailed,
		index.QueueGroup,
		handler.processInvalidNode,
	)
	if err != nil {
		return fmt.Errorf(
			"failed to subscribe to '%s': %v",
			messaging.NodeValidationFailed,
			err,
		)
	}
	return nil
}

// Temp counter for debugging.
var counter uint64

// processValidatedNode handles the processing of validated nodes.
func (handler *nodeHandler) processValidatedNode(msg *natsio.Msg) {
	defer safeAcknowledgeMessage(msg)

	var data messaging.NodeValidatedData
	err := json.Unmarshal(msg.Data, &data)
	if err != nil {
		logger.Error("Failed to unmarshal validated node data", err)
		return
	}

	// Increment the counter
	atomic.AddUint64(&counter, 1)
	logger.Info(fmt.Sprintf("Receiving validated node event no: %d", counter))

	if err = handler.svc.SetNodeValid(&model.Node{
		ProfileURL:  data.ProfileURL,
		ProfileHash: &data.ProfileHash,
		ProfileStr:  data.ProfileStr,
		LastUpdated: &data.LastUpdated,
		Version:     &data.Version,
	}); err != nil {
		logger.Error(
			"Failed to set node valid",
			err,
			zap.String("ProfileURL", data.ProfileURL),
			zap.String("ProfileStr", data.ProfileStr),
		)
	}
}

// processInvalidNode handles the processing of invalid nodes.
func (handler *nodeHandler) processInvalidNode(msg *natsio.Msg) {
	defer safeAcknowledgeMessage(msg)

	var data messaging.NodeValidationFailedData
	err := json.Unmarshal(msg.Data, &data)
	if err != nil {
		logger.Error("Failed to unmarshal invalid node data", err)
		return
	}

	if err = handler.svc.SetNodeInvalid(&model.Node{
		ProfileURL:     data.ProfileURL,
		FailureReasons: data.FailureReasons,
		Version:        &data.Version,
	}); err != nil {
		logger.Error(
			"Failed to set node invalid",
			err,
			zap.String("ProfileURL", data.ProfileURL),
		)
	}
}

// safeAcknowledgeMessage safely acknowledges a message and should be called with
// defer. It recovers from any panics that occurred during message processing and
// then acknowledges the message.
func safeAcknowledgeMessage(msg *natsio.Msg) {
	if err := recover(); err != nil {
		logger.Error(
			"Panic occurred during message processing",
			errors.New("panic"),
			zap.Any("error", err),
		)
	}

	if err := msg.Ack(); err != nil {
		logger.Error("Error acknowledging message", err)
	}
}
