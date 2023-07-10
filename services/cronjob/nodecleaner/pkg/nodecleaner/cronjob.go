package nodecleaner

import (
	"sync"
	"time"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/constant"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/dateutil"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/elastic"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/logger"
	mongodb "github.com/MurmurationsNetwork/MurmurationsServices/pkg/mongo"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/nodecleaner/config"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/nodecleaner/internal/repository/mongo"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/nodecleaner/internal/service"
)

type NodeCleaner struct {
	// Ensures cleanup is only run once.
	runCleanup sync.Once
}

// NewCronJob creates a new instance of NodeCleaner.
func NewCronJob() *NodeCleaner {
	config.Init()

	uri := mongodb.GetURI(
		config.Conf.Mongo.USERNAME,
		config.Conf.Mongo.PASSWORD,
		config.Conf.Mongo.HOST,
	)

	err := mongodb.NewClient(uri, config.Conf.Mongo.DBName)
	if err != nil {
		logger.Panic("Failed to connect to MongoDB", err)
	}

	err = mongodb.Client.Ping()
	if err != nil {
		logger.Panic("Failed to ping MongoDB", err)
	}

	err = elastic.NewClient(config.Conf.ES.URL)
	if err != nil {
		logger.Panic("Failed to connect to Elasticsearch", err)
	}

	return &NodeCleaner{}
}

// Run performs the node cleanup process including the deletion of nodes with
// "validation_failed" status and the deletion of nodes with "deleted" status.
func (nc *NodeCleaner) Run() {
	svc := service.NewNodeService(
		mongo.NewNodeRepository(mongodb.Client.GetClient()),
	)

	err := svc.RemoveValidationFailed()
	if err != nil {
		logger.Panic(
			"Failed to delete nodes with 'validation_failed' status",
			err,
		)
	}

	deletedTimeout := dateutil.NowSubtract(
		time.Duration(config.Conf.TTL.DeletedTTL) * time.Second,
	)

	err = svc.RemoveDeleted(constant.NodeStatus.Deleted, deletedTimeout)
	if err != nil {
		logger.Panic(
			"Failed to delete nodes with 'deleted' status in MongoDB",
			err,
		)
	}

	err = svc.RemoveES(constant.NodeStatus.Deleted, deletedTimeout)
	if err != nil {
		logger.Panic(
			"Failed to delete nodes with 'deleted' status in Elasticsearch",
			err,
		)
	}

	nc.cleanup()
}

// cleanup will clean up the resources associated with the cron job.
func (nc *NodeCleaner) cleanup() {
	nc.runCleanup.Do(func() {
		mongodb.Client.Disconnect()
	})
}
