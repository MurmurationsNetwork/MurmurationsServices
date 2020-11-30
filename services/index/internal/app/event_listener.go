package app

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/common/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/controller/event"
)

func listenToEvents() {
	err := event.HandleNodeValidationFailed().Listen()
	if err != nil {
		logger.Panic("error when trying to listen an event", err)
	}
	err = event.HandleNodeValidated().Listen()
	if err != nil {
		logger.Panic("error when trying to listen an event", err)
	}
}
