package main

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/revalidatenode/pkg/revalidatenode"
)

func main() {
	logger.Info("Starting node revalidation process...")

	s := revalidatenode.NewCronJob()
	if err := s.Run(); err != nil {
		logger.Panic("Failed to revalidate nodes: ", err)
		return
	}

	logger.Info("Node revalidation process completed successfully")
}
