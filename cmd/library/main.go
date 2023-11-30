package main

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/library/pkg/library"
)

func main() {
	logger.Info("Library service starting")

	s := library.NewService()

	go func() {
		<-s.WaitUntilUp()
		logger.Info("Library service started")
	}()

	s.Run()
}
