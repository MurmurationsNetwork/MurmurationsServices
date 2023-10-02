package main

import (
	"context"
	"time"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/nodecleaner/pkg/nodecleaner"
)

func main() {
	logger.Info("Starting NodeCleaner...")

	nc := nodecleaner.NewCronJob()

	startTime := time.Now()

	if err := nc.Run(context.Background()); err != nil {
		logger.Panic("Error running NodeCleaner", err)
		return
	}

	duration := time.Since(startTime)
	logger.Info("NodeCleaner completed successfully")
	logger.Info("NodeCleaner run duration: " + duration.String())
}
