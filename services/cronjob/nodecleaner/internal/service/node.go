package service

import (
	"context"
	"fmt"
	"time"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/constant"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/dateutil"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/nodecleaner/config"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/nodecleaner/internal/repository/es"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/nodecleaner/internal/repository/mongo"
)

// NodesService is an interface that defines methods to remove nodes with specific statuses.
type NodesService interface {
	// Removes nodes with ValidationFailed status.
	RemoveValidationFailed(ctx context.Context) error
	// Removes nodes with Deleted status.
	RemoveDeleted(ctx context.Context) error
}

type nodesService struct {
	mongoRepo mongo.NodeRepository
	esRepo    es.NodeRepository
}

// NewNodeService initializes and returns a new NodesService with the provided
// NodeRepository instances.
func NewNodeService(
	mongoRepo mongo.NodeRepository,
	esRepo es.NodeRepository,
) NodesService {
	return &nodesService{
		mongoRepo: mongoRepo,
		esRepo:    esRepo,
	}
}

// RemoveValidationFailed removes nodes with ValidationFailed status created
// before a calculated time from the MongoDB repository.
func (svc *nodesService) RemoveValidationFailed(ctx context.Context) error {
	timeBefore := dateutil.NowSubtract(
		time.Duration(config.Conf.TTL.ValidationFailedTTL) * time.Second,
	)
	return svc.mongoRepo.RemoveByCreatedAt(
		ctx,
		constant.NodeStatus.ValidationFailed,
		timeBefore,
	)
}

// RemoveDeleted removes nodes with Deleted status updated before the calculated time.
func (svc *nodesService) RemoveDeleted(ctx context.Context) error {
	timeBefore := dateutil.NowSubtract(
		time.Duration(config.Conf.TTL.DeletedTTL) * time.Second,
	)

	err := svc.mongoRepo.RemoveByLastUpdated(
		ctx,
		constant.NodeStatus.Deleted,
		timeBefore,
	)
	if err != nil {
		return fmt.Errorf("error removing nodes from MongoDB: %v", err)
	}

	err = svc.esRepo.Remove(ctx, constant.NodeStatus.Deleted, timeBefore)
	if err != nil {
		return fmt.Errorf("error removing nodes from Elasticsearch: %v", err)
	}

	return nil
}
