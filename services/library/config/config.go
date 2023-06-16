package config

import (
	"time"
)

// Values holds the application configuration.
var Values = Config{}

// Config holds all configuration data for the application.
type Config struct {
	Server serverConfig
	Mongo  mongoConfig
	Static staticConfig
}

// serverConfig holds server-specific configuration.
type serverConfig struct {
	// Server's port.
	Port string `env:"SERVER_PORT,required"`
	// Timeout for reading data.
	TimeoutRead time.Duration `env:"SERVER_TIMEOUT_READ,required"`
	// Timeout for writing data.
	TimeoutWrite time.Duration `env:"SERVER_TIMEOUT_WRITE,required"`
	// Timeout for idle connections.
	TimeoutIdle time.Duration `env:"SERVER_TIMEOUT_IDLE,required"`
	// Rate limit period for GET requests.
	GetRateLimitPeriod string `env:"GET_RATE_LIMIT_PERIOD,required"`
	// Rate limit period for POST requests.
	PostRateLimitPeriod string `env:"POST_RATE_LIMIT_PERIOD,required"`
}

// mongoConfig holds MongoDB-specific configuration.
type mongoConfig struct {
	// MongoDB username.
	Username string `env:"MONGO_USERNAME,required"`
	// MongoDB password.
	Password string `env:"MONGO_PASSWORD,required"`
	// MongoDB host.
	Host string `env:"MONGO_HOST,required"`
	// MongoDB database name.
	DBName string `env:"MONGO_DB_NAME,required"`
}

// staticConfig holds static file-specific configuration.
type staticConfig struct {
	// Path to the static files.
	StaticFilePath string `env:"STATIC_FILE_PATH,required"`
}
