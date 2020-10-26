package services

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/common/constants"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/crypto_utils"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/date_utils"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/events"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/rest_errors"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/datasources/nats"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/domain/nodes"
)

var (
	NodeService nodesServiceInterface = &nodesService{}
)

type nodesServiceInterface interface {
	AddNode(node nodes.Node) (*nodes.Node, rest_errors.RestErr)
	SetNodeValid(node nodes.Node) error
	SetNodeInValid(node nodes.Node) error
	SearchNode(query *nodes.NodeQuery) (nodes.Nodes, rest_errors.RestErr)
	DeleteNode(nodeId string) rest_errors.RestErr
}

type nodesService struct{}

func (s *nodesService) AddNode(node nodes.Node) (*nodes.Node, rest_errors.RestErr) {
	if err := node.Validate(); err != nil {
		return nil, err
	}

	node.NodeID = crypto_utils.GetSHA256(node.ProfileUrl)
	node.Status = constants.Received
	if err := node.Add(); err != nil {
		return nil, err
	}

	events.NewNodeCreatedPublisher(nats.Client()).Publish(events.NodeCreatedData{
		ProfileUrl:    node.ProfileUrl,
		LinkedSchemas: node.LinkedSchemas,
	})

	return &node, nil
}

func (s *nodesService) SetNodeValid(node nodes.Node) error {
	node.NodeID = crypto_utils.GetSHA256(node.ProfileUrl)
	node.Status = constants.Validated
	node.FailedReasons = &[]string{}
	if err := node.Update(); err != nil {
		return err
	}
	return nil
}

func (s *nodesService) SetNodeInValid(node nodes.Node) error {
	node.NodeID = crypto_utils.GetSHA256(node.ProfileUrl)
	node.Status = constants.ValidationFailed
	node.LastValidated = date_utils.GetZeroValueUnix()
	if err := node.Update(); err != nil {
		return err
	}
	return nil
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
