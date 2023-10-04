package revalidatenode

import (
	"sync"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/logger"
	mongodb "github.com/MurmurationsNetwork/MurmurationsServices/pkg/mongo"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/nats"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/revalidatenode/config"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/revalidatenode/internal/repository/mongo"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/revalidatenode/internal/service"
)

// NodeRevalidationCron handles the initialization and running of the node revalidation cron job.
type NodeRevalidationCron struct {
	// Ensures cleanup is run only once.
	runCleanup sync.Once
}

// NewCronJob initializes a new NodeRevalidationCron instance with necessary configurations.
func NewCronJob() *NodeRevalidationCron {
	config.Init()

	// Initialize MongoDB client.
	uri := mongodb.GetURI(
		config.Conf.Mongo.USERNAME,
		config.Conf.Mongo.PASSWORD,
		config.Conf.Mongo.HOST,
	)
	if err := mongodb.NewClient(uri, config.Conf.Mongo.DBName); err != nil {
		logger.Panic("error when trying to connect to MongoDB", err)
	}

	// Check MongoDB connection.
	if err := mongodb.Client.Ping(); err != nil {
		logger.Panic("error when trying to ping the MongoDB", err)
	}

	// Initialize NATS client.
	err := nats.NewClient(
		config.Conf.Nats.ClusterID,
		config.Conf.Nats.ClientID,
		config.Conf.Nats.URL,
	)
	if err != nil {
		logger.Panic("error when trying to connect to NATS", err)
	}

	return &NodeRevalidationCron{}
}

// Run executes the node revalidation process.
func (nc *NodeRevalidationCron) Run() error {
	// Create and run node service for revalidation.
	nodeService := service.NewNodeService(
		mongo.NewNodeRepository(mongodb.Client.GetClient()),
	)

	if err := nodeService.RevalidateNodes(); err != nil {
		return err
	}

	// Perform cleanup after running the service.
	nc.cleanup()

	return nil
}

// cleanup disconnects MongoDB and NATS clients.
func (nc *NodeRevalidationCron) cleanup() {
	nc.runCleanup.Do(func() {
		mongodb.Client.Disconnect()
		nats.Client.Disconnect()
	})
}
