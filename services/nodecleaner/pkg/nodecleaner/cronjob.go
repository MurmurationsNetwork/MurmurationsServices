package nodecleaner

import (
	"context"
	"os"
	"sync"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/elastic"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/logger"
	mongodb "github.com/MurmurationsNetwork/MurmurationsServices/pkg/mongo"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/nodecleaner/config"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/nodecleaner/internal/repository/es"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/nodecleaner/internal/repository/mongo"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/nodecleaner/internal/service"
)

// NodeCleaner manages the node cleanup process.
type NodeCleaner struct {
	runCleanup sync.Once // Ensures cleanup is only run once.
}

// NewCronJob creates a new instance of NodeCleaner.
func NewCronJob() *NodeCleaner {
	config.Init()

	uri := mongodb.GetURI(
		config.Values.Mongo.USERNAME,
		config.Values.Mongo.PASSWORD,
		config.Values.Mongo.HOST,
	)
	if err := mongodb.NewClient(uri, config.Values.Mongo.DBName); err != nil {
		logger.Error("Failed to connect to MongoDB", err)
		os.Exit(1)
	}

	if err := mongodb.Client.Ping(); err != nil {
		logger.Error("Failed to ping MongoDB", err)
		os.Exit(1)
	}

	if err := elastic.NewClient(config.Values.ES.URL); err != nil {
		logger.Error("Failed to connect to Elasticsearch", err)
		os.Exit(1)
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

	// Set nodes with expired status to deleted
	if err := svc.SetExpiredToDeleted(ctx); err != nil {
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
