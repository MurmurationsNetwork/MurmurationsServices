package main

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/schemaparser/pkg/schemaparser"
)

func main() {
	logger.Info("Start loading the schemas from the library repo...")
	s := schemaparser.NewCronJob()
	s.Run()
	logger.Info("Library repo schemas loaded successfully")
}
