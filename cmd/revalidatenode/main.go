package main

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/revalidatenode/pkg/revalidatenode"
)

func main() {
	logger.Info("Start revalidating nodes...")
	s := revalidatenode.NewCronJob()
	s.Run()
	logger.Info("Nodes were successfully revalidated")
}
