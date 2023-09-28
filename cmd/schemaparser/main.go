package main

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/schemaparser/pkg/schemaparser"
)

func main() {
	logger.Info("Start loading schemas...")

	s := schemaparser.NewCronJob()
	if err := s.Run(); err != nil {
		logger.Panic("Error loading schemas", err)
		return
	}

	logger.Info("Schemas were loaded successfully")
}
