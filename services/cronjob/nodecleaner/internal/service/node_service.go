package service

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/nodecleaner/internal/repository/db"
)

type NodesService interface {
	Remove(status string, timeBefore int64) error
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

func (svc *nodesService) Remove(status string, timeBefore int64) error {
	return svc.nodeRepo.Remove(status, timeBefore)
}

func (svc *nodesService) RemoveES(status string, timeBefore int64) error {
	return svc.nodeRepo.RemoveES(status, timeBefore)
}
