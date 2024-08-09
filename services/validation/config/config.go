package config

import (
	"time"
)

// Values holds the application configuration.
var Values = Config{}

// Config represents the application configuration structure.
type Config struct {
	Server  ServerConfig
	Library LibraryConfig
	NATS    NATSConfig
	Redis   redisConf
}

// ServerConfig holds the server related configuration.
type ServerConfig struct {
	// Server port
	Port string `env:"SERVER_PORT,required"`
	// Server read timeout
	TimeoutRead time.Duration `env:"SERVER_TIMEOUT_READ,required"`
	// Server write timeout
	TimeoutWrite time.Duration `env:"SERVER_TIMEOUT_WRITE,required"`
	// Server idle timeout
	TimeoutIdle time.Duration `env:"SERVER_TIMEOUT_IDLE,required"`
}

// LibraryConfig holds the library related configuration.
type LibraryConfig struct {
	// Internal URL of the library
	InternalURL string `env:"LIBRARY_URL,required"`
}

// NATSConfig holds the NATS related configuration.
type NATSConfig struct {
	// NATS cluster ID
	ClusterID string `env:"NATS_CLUSTER_ID,required"`
	// NATS client ID
	ClientID string `env:"NATS_CLIENT_ID,required"`
	// NATS URL
	URL string `env:"NATS_URL,required"`
}

type redisConf struct {
	URL string `env:"REDIS_URL,required"`
}
