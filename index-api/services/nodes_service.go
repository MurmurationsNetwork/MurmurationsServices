package services

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/index-api/domain/nodes"
	"github.com/MurmurationsNetwork/MurmurationsServices/utils/rest_errors"
)

var (
	NodeService nodesServiceInterface = &nodesService{}
)

type nodesServiceInterface interface {
	AddNode(node nodes.Node) (*nodes.Node, rest_errors.RestErr)
	GetNode(nodeId string) (*nodes.Node, rest_errors.RestErr)
}

type nodesService struct{}

func (s *nodesService) AddNode(node nodes.Node) (*nodes.Node, rest_errors.RestErr) {
	if err := node.Validate(); err != nil {
		return nil, err
	}
	if err := node.Add(); err != nil {
		return nil, err
	}
	return &node, nil
}

func (s *nodesService) GetNode(nodeId string) (*nodes.Node, rest_errors.RestErr) {
	result := nodes.Node{ID: nodeId}
	if err := result.Get(); err != nil {
		return nil, err
	}
	return &result, nil
}
