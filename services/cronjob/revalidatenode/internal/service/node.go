package service

import (
	"fmt"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/constant"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/event"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/nats"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/revalidatenode/internal/repository/mongo"
)

type NodeService interface {
	RevalidateNodes() error
}

type nodeService struct {
	mongoRepo mongo.NodeRepository
}

func NewNodeService(mongoRepo mongo.NodeRepository) NodeService {
	return &nodeService{
		mongoRepo: mongoRepo,
	}
}

func (svc *nodeService) RevalidateNodes() error {
	statuses := []string{
		constant.NodeStatus.Received,
		constant.NodeStatus.PostFailed,
	}

	nodes, err := svc.mongoRepo.FindByStatuses(statuses)
	if err != nil {
		return err
	}

	if len(nodes) == 0 {
		return nil
	}

	logger.Info(
		fmt.Sprintf(
			"Found %d nodes with status %s or %s, sending them to validation service",
			len(nodes),
			statuses[0],
			statuses[1],
		),
	)

	for _, node := range nodes {
		event.NewNodeCreatedPublisher(nats.Client.Client()).
			Publish(event.NodeCreatedData{
				ProfileURL: node.ProfileURL,
				Version:    *node.Version,
			})
	}

	return nil
}
