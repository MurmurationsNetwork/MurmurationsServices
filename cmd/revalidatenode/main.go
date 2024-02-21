package main

import (
	"os"
	"time"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/revalidatenode/pkg/revalidatenode"
)

func main() {
	s := revalidatenode.NewCronJob()

	startTime := time.Now()

	if err := s.Run(); err != nil {
		logger.Error("Failed to revalidate nodes: ", err)
		os.Exit(1)
	}

	duration := time.Since(startTime)
	logger.Info("Node revalidation run duration: " + duration.String())
}
