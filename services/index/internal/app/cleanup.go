package app

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/mongo"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/nats"
)

func cleanup() {
	mongo.Client.Disconnect()
	nats.Client.Disconnect()
}
