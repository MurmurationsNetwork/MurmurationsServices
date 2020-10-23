package app

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/common/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/validation/internal/events"
)

func listenToEvents() {
	err := events.HandleNodeCreated.Listen()
	if err != nil {
		logger.Panic("error when trying to listen an event", err)
	}
}
