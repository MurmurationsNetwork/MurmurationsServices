package service

import (
	"time"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/constant"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/dateutil"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/nodecleaner/config"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/nodecleaner/internal/repository/db"
)

type NodesService interface {
	RemoveValidationFailed() error
	RemoveDeleted(status string, timeBefore int64) error
	RemoveES(status string, timeBefore int64) error
}

type nodesService struct {
	nodeRepo db.NodeRepository
}

func NewNodeService(nodeRepo db.NodeRepository) NodesService {
	return &nodesService{
		nodeRepo: nodeRepo,
	}
}

func (svc *nodesService) RemoveValidationFailed() error {
	return svc.nodeRepo.RemoveValidationFailed(
		constant.NodeStatus.ValidationFailed,
		dateutil.NowSubtract(
			time.Duration(config.Conf.TTL.ValidationFailedTTL)*time.Second,
		),
	)
}

func (svc *nodesService) RemoveDeleted(status string, timeBefore int64) error {
	return svc.nodeRepo.RemoveDeleted(status, timeBefore)
}

func (svc *nodesService) RemoveES(status string, timeBefore int64) error {
	return svc.nodeRepo.RemoveES(status, timeBefore)
}
