package usecase

import (
	"fmt"

	"github.com/MurmurationsNetwork/MurmurationsServices/common/constant"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/event"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/nats"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/revalidatenode/internal/adapter/repository"
)

type NodeUsecase interface {
	RevalidateNodes() error
}

type nodeUsecase struct {
	nodeRepo repository.NodeRepository
}

func NewNodeUsecase(nodeRepo repository.NodeRepository) NodeUsecase {
	return &nodeUsecase{
		nodeRepo: nodeRepo,
	}
}

func (usecase *nodeUsecase) RevalidateNodes() error {
	nodes, err := usecase.nodeRepo.FindByStatus(constant.NodeStatus.Received)
	if err != nil {
		return err
	}

	if len(nodes) == 0 {
		return nil
	}

	logger.Info(fmt.Sprintf("Found %d nodes with status %s, sending them to validation service", len(nodes), constant.NodeStatus.Received))

	for _, node := range nodes {
		event.NewNodeCreatedPublisher(nats.Client.Client()).Publish(event.NodeCreatedData{
			ProfileURL: node.ProfileURL,
			Version:    *node.Version,
		})
	}

	return nil
}
