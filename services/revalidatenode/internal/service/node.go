package service

import (
	"context"
	"fmt"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/constant"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/event"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/nats"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/revalidatenode/internal/repository/mongo"
)

// NodeService outlines methods to interact with node data.
type NodeService interface {
	RevalidateNodes() error
}

// nodeService implements NodeService with a mongo.NodeRepository.
type nodeService struct {
	mongoRepo mongo.NodeRepository
}

// NewNodeService initializes a new nodeService instance.
func NewNodeService(mongoRepo mongo.NodeRepository) NodeService {
	return &nodeService{mongoRepo: mongoRepo}
}

// RevalidateNodes processes nodes with specific statuses and sends them for re-validation.
func (svc *nodeService) RevalidateNodes() error {
	statuses := []string{
		constant.NodeStatus.Received,
		constant.NodeStatus.PostFailed,
	}

	page := 1
	pageSize := 100

	for {
		nodes, err := svc.mongoRepo.FindByStatuses(
			context.Background(),
			statuses,
			page,
			pageSize,
		)
		if err != nil {
			return err
		}

		if len(nodes) == 0 {
			break
		}

		logger.Info(
			fmt.Sprintf(
				"Found %d nodes with status %s or %s on page %d, sending them to validation service",
				len(nodes),
				statuses[0],
				statuses[1],
				page,
			),
		)

		for _, node := range nodes {
			err := event.NewNodeCreatedPublisher(nats.Client.JetStream()).
				PublishSync(event.NodeCreatedData{
					ProfileURL: node.ProfileURL,
					Version:    *node.Version,
				})
			if err != nil {
				logger.Error("Failed to publish node:created event: ", err)
			}
		}

		if len(nodes) < pageSize {
			break
		}

		page++
	}

	return nil
}
