package app

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/services/validation/internal/controller/event"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/validation/internal/service"
)

func listenToEvents() error {
	nodeHandler := event.NewNodeHandler(service.NewValidationService())
	err := nodeHandler.NewNodeCreatedListener()
	if err != nil {
		return err
	}
	return nil
}
