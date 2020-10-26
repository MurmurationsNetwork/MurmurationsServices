package app

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/common/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/validation/internal/controller/event"
)

func listenToEvents() {
	err := event.HandleNodeCreated.Listen()
	if err != nil {
		logger.Panic("error when trying to listen an event", err)
	}
}
