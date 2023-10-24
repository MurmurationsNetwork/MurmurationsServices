package main

import (
	"os"
	"time"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/dataproxyrefresher/pkg/fresher"
)

func main() {
	logger.Info("Start DataproxyRefresher...")
	startTime := time.Now()

	refresher := fresher.NewRefresher()

	if err := refresher.Run(); err != nil {
		logger.Error("Error running DataproxyRefresher", err)
		os.Exit(1)
		return
	}

	duration := time.Since(startTime)
	logger.Info("DataproxyRefresher has finished")
	logger.Info("DataproxyRefresher run duration: " + duration.String())
}
