package event

import (
	"encoding/json"
	"errors"
	"fmt"

	stan "github.com/nats-io/stan.go"

	"github.com/MurmurationsNetwork/MurmurationsServices/common/event"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/nats"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/validation/internal/entity"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/validation/internal/service"
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
		// If we don't put a goruine here, we can only validate a node each time.
		go func() {
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

			handler.validationService.ValidateNode(&entity.Node{
				ProfileURL: nodeCreatedData.ProfileURL,
				Version:    nodeCreatedData.Version,
			})

			_ = msg.Ack()
		}()
	}).Listen()
}
