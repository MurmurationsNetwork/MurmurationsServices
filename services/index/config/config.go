// Package config manages all configurations in the application
package config

import (
	"time"
)

// Values holds all the application's configurations.
var Values = config{
	FeatureToggles: make(map[string]bool),
}

// config is the central configuration holder.
type config struct {
	// Server configuration
	Server serverConf
	// Library configuration
	Library libraryConf
	// MongoDB configuration
	Mongo mongoConf
	// Elasticsearch configuration
	ES esConf
	// NATS configuration
	Nats natsConf
	// TTL configuration
	TTL ttlConf
	// FeatureToggles
	FeatureToggles map[string]bool
}

// serverConf contains server configurations.
type serverConf struct {
	// Server port
	Port string `env:"SERVER_PORT,required"`
	// Server read timeout
	TimeoutRead time.Duration `env:"SERVER_TIMEOUT_READ,required"`
	// Server write timeout
	TimeoutWrite time.Duration `env:"SERVER_TIMEOUT_WRITE,required"`
	// Server idle timeout
	TimeoutIdle time.Duration `env:"SERVER_TIMEOUT_IDLE,required"`
	// Rate limit period for GET requests
	GetRateLimitPeriod string `env:"GET_RATE_LIMIT_PERIOD,required"`
	// Rate limit period for POST requests
	PostRateLimitPeriod string `env:"POST_RATE_LIMIT_PERIOD,required"`
	// Maximum number of tags in an array
	TagsArraySize string `env:"TAGS_ARRAY_SIZE,required"`
	// Maximum length of a tag string
	TagsStringLength string `env:"TAGS_STRING_LENGTH,required"`
	// Fuzziness of tag matching
	TagsFuzziness string `env:"TAGS_FUZZINESS,required"`
}

// libraryConf contains configuration for the internal library.
type libraryConf struct {
	// URL for the internal library
	InternalURL string `env:"LIBRARY_URL,required"`
}

// mongoConf contains MongoDB configuration details.
type mongoConf struct {
	// MongoDB username
	USERNAME string `env:"MONGO_USERNAME,required"`
	// MongoDB password
	PASSWORD string `env:"MONGO_PASSWORD,required"`
	// MongoDB host address
	HOST string `env:"MONGO_HOST,required"`
	// MongoDB database name
	DBName string `env:"MONGO_DB_NAME,required"`
}

// esConf contains the configuration for the Elasticsearch service.
type esConf struct {
	// Elasticsearch service URL
	URL string `env:"ELASTICSEARCH_URL,required"`
}

// natsConf contains the configuration for the NATS service.
type natsConf struct {
	// NATS cluster ID
	ClusterID string `env:"NATS_CLUSTER_ID,required"`
	// NATS client ID
	ClientID string `env:"NATS_CLIENT_ID,required"`
	// NATS service URL
	URL string `env:"NATS_URL,required"`
}

// ttlConf contains the configuration for the TTL (Time To Live) settings.
type ttlConf struct {
	// Time To Live for deleted items.
	DeletedTTL int64 `env:"DELETED_TTL,required"`
}
