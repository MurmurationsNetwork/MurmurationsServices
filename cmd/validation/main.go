package main

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/validation/pkg/validation"
)

func main() {
	logger.Info("Validation service starting")

	s := validation.NewService()

	go func() {
		<-s.WaitUntilUp()
		logger.Info("Validation service started")
	}()

	s.Run()
}
