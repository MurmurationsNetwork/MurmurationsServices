package services

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/indexer/internal/domain/nodes"
	"github.com/MurmurationsNetwork/MurmurationsServices/utils/date_utils"
	"github.com/MurmurationsNetwork/MurmurationsServices/utils/hash_utils"
	"github.com/MurmurationsNetwork/MurmurationsServices/utils/rest_errors"
)

var (
	NodeService nodesServiceInterface = &nodesService{}
)

type nodesServiceInterface interface {
	AddNode(node nodes.Node) (*nodes.Node, rest_errors.RestErr)
	GetNode(nodeId string) (*nodes.Node, rest_errors.RestErr)
	SearchNode(query *nodes.NodeQuery) (nodes.Nodes, rest_errors.RestErr)
	DeleteNode(nodeId string) rest_errors.RestErr
}

type nodesService struct{}

func (s *nodesService) AddNode(node nodes.Node) (*nodes.Node, rest_errors.RestErr) {
	if err := node.Validate(); err != nil {
		return nil, err
	}

	node.NodeID = hash_utils.SHA256(node.ProfileUrl)
	node.LastValidated = date_utils.GetNowUnix()

	if err := node.Add(); err != nil {
		return nil, err
	}
	return &node, nil
}

func (s *nodesService) GetNode(nodeId string) (*nodes.Node, rest_errors.RestErr) {
	node := &nodes.Node{NodeID: nodeId}
	if err := node.Get(); err != nil {
		return nil, err
	}
	return node, nil
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
