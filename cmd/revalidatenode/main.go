package main

import (
	"time"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/revalidatenode/pkg/revalidatenode"
)

func main() {
	logger.Info("Starting node revalidation process...")

	s := revalidatenode.NewCronJob()

	startTime := time.Now()

	if err := s.Run(); err != nil {
		logger.Panic("Failed to revalidate nodes: ", err)
		return
	}

	duration := time.Since(startTime)
	logger.Info("Node revalidation process completed successfully")
	logger.Info("Node revalidation run duration: " + duration.String())
}
