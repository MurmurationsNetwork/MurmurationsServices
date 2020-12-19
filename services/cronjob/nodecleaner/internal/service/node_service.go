package service

import (
	"time"

	"github.com/MurmurationsNetwork/MurmurationsServices/common/constant"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/dateutil"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/nodecleaner/internal/config"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/nodecleaner/internal/repository/db"
)

type NodesService interface {
	Remove() error
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
		dateutil.NowSubtract(time.Duration(config.Conf.TTL)*time.Second),
	)
}
