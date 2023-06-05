package main

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/schemaparser/pkg/schemaparser"
)

func main() {
	logger.Info("Start loading schemas...")
	s := schemaparser.NewCronJob()
	s.Run()
	logger.Info("Schemas were loaded successfully")
}
