package dataproxy

import (
	"log"

	"github.com/MurmurationsNetwork/MurmurationsServices/services/dataproxy/config"
	env "github.com/caarlos0/env/v10"
)

func init() {
	err := env.Parse(&config.Values)
	if err != nil {
		log.Fatalf("Failed to decode environment variables: %s", err)
	}
}
