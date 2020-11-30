package app

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/common/nats"
)

func cleanup() {
	nats.Client.Disconnect()
}
