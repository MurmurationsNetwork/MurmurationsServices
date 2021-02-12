package main

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/common/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/mongo"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/nodecleaner/global"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/nodecleaner/internal/repository/db"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/nodecleaner/internal/service"
)

func init() {
	global.Init()
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
