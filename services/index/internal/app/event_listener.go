package app

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/controller/event"
)

func listenToEvents() error {
	err := event.HandleNodeValidationFailed().Listen()
	if err != nil {
		return err
	}
	err = event.HandleNodeValidated().Listen()
	if err != nil {
		return err
	}
	return nil
}
