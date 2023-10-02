package nodecleaner

import (
	"context"
	"sync"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/elastic"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/logger"
	mongodb "github.com/MurmurationsNetwork/MurmurationsServices/pkg/mongo"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/nodecleaner/config"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/nodecleaner/internal/repository/es"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/nodecleaner/internal/repository/mongo"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/nodecleaner/internal/service"
)

// NodeCleaner manages the node cleanup process.
type NodeCleaner struct {
	runCleanup sync.Once // Ensures cleanup is only run once.
}

// NewCronJob creates a new instance of NodeCleaner.
func NewCronJob() *NodeCleaner {
	config.Init()

	uri := mongodb.GetURI(
		config.Conf.Mongo.USERNAME,
		config.Conf.Mongo.PASSWORD,
		config.Conf.Mongo.HOST,
	)
	if err := mongodb.NewClient(uri, config.Conf.Mongo.DBName); err != nil {
		logger.Panic("Failed to connect to MongoDB", err)
	}

	if err := mongodb.Client.Ping(); err != nil {
		logger.Panic("Failed to ping MongoDB", err)
	}

	if err := elastic.NewClient(config.Conf.ES.URL); err != nil {
		logger.Panic("Failed to connect to Elasticsearch", err)
	}

	return &NodeCleaner{}
}

// Run executes the node cleanup process.
func (nc *NodeCleaner) Run(ctx context.Context) error {
	svc := service.NewNodeService(
		mongo.NewNodeRepository(mongodb.Client.GetClient()),
		es.NewNodeRepository(),
	)

	if err := svc.RemoveValidationFailed(ctx); err != nil {
		return err
	}

	if err := svc.RemoveDeleted(ctx); err != nil {
		return err
	}

	nc.cleanup()
	return nil
}

// cleanup releases resources associated with the NodeCleaner.
func (nc *NodeCleaner) cleanup() {
	nc.runCleanup.Do(func() {
		mongodb.Client.Disconnect()
	})
}
