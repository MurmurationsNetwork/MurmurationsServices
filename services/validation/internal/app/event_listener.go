package app

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/services/validation/internal/controller/event"
)

func listenToEvents() error {
	err := event.HandleNodeCreated().Listen()
	if err != nil {
		return err
	}
	return nil
}
