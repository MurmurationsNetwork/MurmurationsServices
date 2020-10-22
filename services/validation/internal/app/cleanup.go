package app

import "github.com/MurmurationsNetwork/MurmurationsServices/services/validation/internal/datasources/nats"

func cleanup() {
	nats.Disconnect()
}
