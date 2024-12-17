package config

import (
	"log"

	env "github.com/caarlos0/env/v10"
)

// Conf holds the configuration settings for the application.
var Values = config{}

type config struct {
	// Mongo holds the configuration for MongoDB.
	Mongo mongoConf
	// Nats holds the configuration for NATS.
	Nats natsConf
}

type mongoConf struct {
	// USERNAME is the user name used to authenticate with MongoDB.
	USERNAME string `env:"MONGO_USERNAME,required"`
	// PASSWORD is the password used to authenticate with MongoDB.
	PASSWORD string `env:"MONGO_PASSWORD,required"`
	// HOST is the host address for MongoDB.
	HOST string `env:"MONGO_HOST,required"`
	// DBName is the name of the MongoDB database to connect to.
	DBName string `env:"MONGO_DB_NAME,required"`
}

type natsConf struct {
	// ClusterID is the NATS cluster identifier.
	ClusterID string `env:"NATS_CLUSTER_ID,required"`
	// ClientID is the NATS client identifier.
	ClientID string `env:"NATS_CLIENT_ID,required"`
	// URL is the URL used to connect to the NATS server.
	URL string `env:"NATS_URL,required"`
}

// Init initializes the Conf variable by parsing environment variables.
func Init() {
	if err := env.Parse(&Values); err != nil {
		log.Fatalf("Failed to decode environment variables: %s", err)
	}
}
