package schemaparser

import (
	"log"

	env "github.com/caarlos0/env/v6"

	"github.com/MurmurationsNetwork/MurmurationsServices/services/schemaparser/config"
)

func init() {
	err := env.Parse(&config.Values)
	if err != nil {
		log.Fatalf("Failed to decode environment variables: %s", err)
	}
}
