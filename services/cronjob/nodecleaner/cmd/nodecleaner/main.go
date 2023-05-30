package main

import (
	"time"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/constant"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/dateutil"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/mongo"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/nodecleaner/config"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/nodecleaner/global"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/nodecleaner/internal/repository/db"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/nodecleaner/internal/service"
)

func init() {
	global.Init()
}

func main() {
	svc := service.NewNodeService(
		db.NewNodeRepository(mongo.Client.GetClient()),
	)

	err := svc.Remove()
	if err != nil {
		logger.Panic(
			"Error when trying to delete nodes with validation_failed status",
			err,
		)
		return
	}

	// delete mongoDB data
	deletedTimeout := dateutil.NowSubtract(
		time.Duration(config.Conf.TTL.DeletedTTL) * time.Second,
	)

	err = svc.RemoveDeleted(constant.NodeStatus.Deleted, deletedTimeout)
	if err != nil {
		logger.Panic(
			"Error when trying to delete nodes with deleted status",
			err,
		)
		return
	}

	mongo.Client.Disconnect()

	// delete ElasticSearch data
	err = svc.RemoveES(constant.NodeStatus.Deleted, deletedTimeout)
	if err != nil {
		logger.Panic(
			"Error when trying to delete nodes with deleted status",
			err,
		)
		return
	}
}
