package redisadapter

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/common/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/redis"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/schemaparser/config"
)

func NewClient() redis.Redis {
	client := redis.NewClient(config.Conf.Redis.URL)
	err := client.Ping()
	if err != nil {
		logger.Panic("error when trying to ping the redis", err)
		return nil
	}
	return client
}
