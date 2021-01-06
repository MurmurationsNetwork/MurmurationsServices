package main

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/common/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/mongo"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/nodecleaner/config"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/nodecleaner/internal/adapter/mongodb"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/nodecleaner/internal/repository/db"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/nodecleaner/internal/service"
)

func init() {
	config.Init()
	mongodb.Init()
}

func main() {
	svc := service.NewNodeService(db.NewNodeRepository(mongo.Client.GetClient()))

	err := svc.Remove()
	if err != nil {
		logger.Panic("Error when trying to delete nodes", err)
		return
	}

	mongo.Client.Disconnect()
}
