package library

import (
	"log"

	env "github.com/caarlos0/env/v10"

	"github.com/MurmurationsNetwork/MurmurationsServices/services/library/config"
)

func init() {
	err := env.Parse(&config.Values)
	if err != nil {
		log.Fatalf("Failed to decode environment variables: %s", err)
	}
}
