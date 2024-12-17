package revalidatenode

import (
	"os"
	"sync"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/logger"
	mongodb "github.com/MurmurationsNetwork/MurmurationsServices/pkg/mongo"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/natsclient"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/revalidatenode/config"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/revalidatenode/internal/repository/mongo"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/revalidatenode/internal/service"
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
		config.Values.Mongo.USERNAME,
		config.Values.Mongo.PASSWORD,
		config.Values.Mongo.HOST,
	)
	if err := mongodb.NewClient(uri, config.Values.Mongo.DBName); err != nil {
		logger.Error("error when trying to connect to MongoDB", err)
		os.Exit(1)
	}

	// Check MongoDB connection.
	if err := mongodb.Client.Ping(); err != nil {
		logger.Error("error when trying to ping the MongoDB", err)
		os.Exit(1)
	}

	// Initialize NATS client.
	setupNATS()

	return &NodeRevalidationCron{}
}

// setupNATS initializes Nats service.
func setupNATS() {
	err := natsclient.Initialize(config.Values.Nats.URL)
	if err != nil {
		logger.Error("Failed to create Nats client", err)
		os.Exit(1)
	}
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
		// Disconnect from MongoDB.
		mongodb.Client.Disconnect()

		// Disconnect from NATS.
		if err := natsclient.GetInstance().Disconnect(); err != nil {
			logger.Error("Error disconnecting from NATS: %v", err)
		}
	})
}
