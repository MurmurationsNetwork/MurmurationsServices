package app

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/common/mongo"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/nats"
)

func cleanup() {
	mongo.Client.Disconnect()
	nats.Client.Disconnect()
}
