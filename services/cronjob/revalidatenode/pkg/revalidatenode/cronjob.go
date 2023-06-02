package revalidatenode

import (
	"sync"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/mongo"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/nats"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/revalidatenode/config"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/revalidatenode/internal/adapter/repository"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/revalidatenode/internal/usecase"
)

type NodeRevalidationCron struct {
	// Ensures cleanup is only run once
	runCleanup sync.Once
}

// NewCronJob creates a new instance of NodeRevalidationCron.
func NewCronJob() *NodeRevalidationCron {
	config.Init()

	uri := mongo.GetURI(
		config.Conf.Mongo.USERNAME,
		config.Conf.Mongo.PASSWORD,
		config.Conf.Mongo.HOST,
	)

	err := mongo.NewClient(uri, config.Conf.Mongo.DBName)
	if err != nil {
		logger.Panic("error when trying to connect to MongoDB", err)
	}

	err = mongo.Client.Ping()
	if err != nil {
		logger.Panic("error when trying to ping the MongoDB", err)
	}

	err = nats.NewClient(
		config.Conf.Nats.ClusterID,
		config.Conf.Nats.ClientID,
		config.Conf.Nats.URL,
	)
	if err != nil {
		logger.Panic("error when trying to connect to NATS", err)
	}

	return &NodeRevalidationCron{}
}

// Run revalidates the nodes.
func (nc *NodeRevalidationCron) Run() {
	nodeUsecase := usecase.NewNodeUsecase(
		(repository.NewNodeRepository(mongo.Client.GetClient())),
	)

	err := nodeUsecase.RevalidateNodes()
	if err != nil {
		logger.Panic("error when revalidating nodes", err)
	}

	nc.cleanup()
}

// cleanup will clean up the resources associated with the cron job.
func (nc *NodeRevalidationCron) cleanup() {
	nc.runCleanup.Do(func() {
		mongo.Client.Disconnect()
		nats.Client.Disconnect()
	})
}
