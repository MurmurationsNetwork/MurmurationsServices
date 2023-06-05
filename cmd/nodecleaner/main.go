package main

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/nodecleaner/pkg/nodecleaner"
)

func main() {
	logger.Info("Start cleaning up nodes...")
	s := nodecleaner.NewCronJob()
	s.Run()
	logger.Info("Nodes were cleaned up successfully")
}
