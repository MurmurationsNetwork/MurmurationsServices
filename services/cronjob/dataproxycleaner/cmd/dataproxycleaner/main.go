package main

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/common/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/mongo"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/dataproxycleaner/config"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/dataproxycleaner/global"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/dataproxycleaner/internal/repository/db"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/dataproxycleaner/internal/service"
	"os"
	"time"
)

func init() {
	global.Init()
}

func cleanUp() {
	mongo.Client.Disconnect()
	os.Exit(0)
}

func main() {
	svc := service.NewProfileService(db.NewProfileRepository(mongo.Client.GetClient()))

	curTime := time.Now().Unix()
	deleteBefore := curTime - config.Conf.DeletedTTL
	profiles, err := svc.FindLessThan(deleteBefore)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			logger.Error("no profile found", err)
		} else {
			logger.Error("failed to find data from profiles", err)
		}
		cleanUp()
	}

	// delete from profile
	for _, profile := range profiles {
		err := svc.Delete(profile.Cuid)
		if err != nil {
			logger.Error("failed to delete data from profiles, profile cuid:"+profile.Cuid, err)
			cleanUp()
		}
	}

	// todo: send delete to index
	cleanUp()
}
