package app

import "github.com/MurmurationsNetwork/MurmurationsServices/services/validation/internal/datasource/nats"

func cleanup() {
	nats.Disconnect()
}
