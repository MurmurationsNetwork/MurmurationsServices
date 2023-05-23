package service

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/common/constant"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/dateutil"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/nodecleaner/config"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/nodecleaner/internal/repository/db"
	"time"
)

type NodesService interface {
	Remove() error
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

func (svc *nodesService) Remove() error {
	return svc.nodeRepo.Remove(
		constant.NodeStatus.ValidationFailed,
		dateutil.NowSubtract(time.Duration(config.Conf.TTL.TTL)*time.Second),
	)
}

func (svc *nodesService) RemoveDeleted(status string, timeBefore int64) error {
	return svc.nodeRepo.RemoveDeleted(status, timeBefore)
}

func (svc *nodesService) RemoveES(status string, timeBefore int64) error {
	return svc.nodeRepo.RemoveES(status, timeBefore)
}
