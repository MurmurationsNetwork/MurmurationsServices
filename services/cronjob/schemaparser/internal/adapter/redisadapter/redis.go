package redisadapter

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/redis"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/schemaparser/config"
)

func NewClient() redis.Redis {
	client := redis.NewClient(config.Conf.Redis.URL)
	err := client.Ping()
	if err != nil {
		logger.Panic("error when trying to ping Redis", err)
		return nil
	}
	return client
}
