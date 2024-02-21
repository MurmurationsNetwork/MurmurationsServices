package main

import (
	"time"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/schemaparser/pkg/schemaparser"
)

func main() {
	startTime := time.Now()

	s := schemaparser.NewCronJob()
	if err := s.Run(); err != nil {
		logger.Panic("Error running SchemaParser", err)
		return
	}

	// Calculate and log the duration
	duration := time.Since(startTime)
	logger.Info("SchemaParser run duration: " + duration.String())
}
