package main

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/common/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/mongo"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/dataproxyupdater/global"
)

func init() {
	global.Init()
}

func main() {
	logger.Info("dataproxyupdater is running!")

	mongo.Client.Disconnect()
}
