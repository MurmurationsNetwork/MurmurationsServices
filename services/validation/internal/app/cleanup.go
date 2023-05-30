package app

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/nats"
)

func cleanup() {
	nats.Client.Disconnect()
}
