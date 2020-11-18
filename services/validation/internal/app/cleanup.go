package app

import "github.com/MurmurationsNetwork/MurmurationsServices/services/validation/internal/adapter/nats"

func cleanup() {
	nats.Disconnect()
}
