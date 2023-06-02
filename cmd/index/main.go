package main

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/pkg/index"
)

func main() {
	logger.Info("Index service starting")

	s := index.NewService()

	go func() {
		<-s.WaitUntilUp()
		logger.Info("Index service started")
	}()

	s.Run()
}
