package service

import (
	"time"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/constant"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/dateutil"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/nodecleaner/config"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/nodecleaner/internal/repository/mongo"
)

type NodesService interface {
	RemoveValidationFailed() error
	RemoveDeleted(status string, timeBefore int64) error
	RemoveES(status string, timeBefore int64) error
}

type nodesService struct {
	mongoRepo mongo.NodeRepository
}

func NewNodeService(mongoRepo mongo.NodeRepository) NodesService {
	return &nodesService{
		mongoRepo: mongoRepo,
	}
}

func (svc *nodesService) RemoveValidationFailed() error {
	return svc.mongoRepo.RemoveValidationFailed(
		constant.NodeStatus.ValidationFailed,
		dateutil.NowSubtract(
			time.Duration(config.Conf.TTL.ValidationFailedTTL)*time.Second,
		),
	)
}

func (svc *nodesService) RemoveDeleted(status string, timeBefore int64) error {
	return svc.mongoRepo.RemoveDeleted(status, timeBefore)
}

func (svc *nodesService) RemoveES(status string, timeBefore int64) error {
	return svc.mongoRepo.RemoveES(status, timeBefore)
}
