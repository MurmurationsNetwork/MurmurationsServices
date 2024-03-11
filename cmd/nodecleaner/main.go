package main

import (
	"context"
	"os"
	"time"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/nodecleaner/pkg/nodecleaner"
)

func main() {
	nc := nodecleaner.NewCronJob()

	startTime := time.Now()

	if err := nc.Run(context.Background()); err != nil {
		logger.Error("Error running NodeCleaner", err)
		os.Exit(1)
	}

	duration := time.Since(startTime)
	logger.Info("NodeCleaner run duration: " + duration.String())
}
