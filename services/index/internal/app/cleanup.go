package app

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/common/mongo"
)

func cleanup() {
	mongo.Client.Disconnect()
}
