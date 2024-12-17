package main

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/dataproxy/pkg/dataproxy"
)

func main() {
	logger.Info("Dataproxy service starting")

	s := dataproxy.NewService()

	go func() {
		<-s.WaitUntilUp()
		logger.Info("Dataproxy service started")
	}()

	s.Run()
}
