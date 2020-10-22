package services

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/common/crypto_utils"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/date_utils"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/events"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/http_utils"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/rest_errors"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/domain/nodes"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/events/publishers"
)

var (
	NodeService nodesServiceInterface = &nodesService{}
)

type nodesServiceInterface interface {
	AddNode(node nodes.Node) (*nodes.Node, rest_errors.RestErr)
	SearchNode(query *nodes.NodeQuery) (nodes.Nodes, rest_errors.RestErr)
	DeleteNode(nodeId string) rest_errors.RestErr
}

type nodesService struct{}

func (s *nodesService) AddNode(node nodes.Node) (*nodes.Node, rest_errors.RestErr) {
	if err := node.Validate(); err != nil {
		return nil, err
	}

	jsonStr, err := http_utils.GetStr(node.ProfileUrl)
	if err != nil {
		return nil, rest_errors.NewBadRequestError(err.Error())
	}
	node.NodeID = crypto_utils.GetSHA256(jsonStr)
	node.LastValidated = date_utils.GetNowUnix()

	if err := node.Add(); err != nil {
		return nil, err
	}

	publishErr := publishers.NodeCreated.Publish(&events.NodeCreatedData{
		ProfileUrl:    node.ProfileUrl,
		LinkedSchemas: node.LinkedSchemas,
	})
	logger.Error("error when trying to publish the node:created event", publishErr)

	return &node, nil
}

func (s *nodesService) SearchNode(query *nodes.NodeQuery) (nodes.Nodes, rest_errors.RestErr) {
	node := &nodes.Node{}
	result, err := node.Search(query)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *nodesService) DeleteNode(nodeId string) rest_errors.RestErr {
	node := &nodes.Node{NodeID: nodeId}
	if err := node.Delete(); err != nil {
		return err
	}
	return nil
}
