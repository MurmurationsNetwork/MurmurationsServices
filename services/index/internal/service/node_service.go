package service

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/common/constant"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/cryptoutil"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/dateutil"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/event"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/resterr"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/datasources/nats"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/domain/node"
)

var (
	NodeService nodesServiceInterface = &nodesService{}
)

type nodesServiceInterface interface {
	AddNode(node node.Node) (*node.Node, resterr.RestErr)
	SetNodeValid(node node.Node) error
	SetNodeInValid(node node.Node) error
	SearchNode(query *node.NodeQuery) (node.Nodes, resterr.RestErr)
	DeleteNode(nodeId string) resterr.RestErr
}

type nodesService struct{}

func (s *nodesService) AddNode(node node.Node) (*node.Node, resterr.RestErr) {
	if err := node.Validate(); err != nil {
		return nil, err
	}

	node.ID = cryptoutil.GetSHA256(node.ProfileUrl)
	node.Status = constant.Received
	if err := node.Add(); err != nil {
		return nil, err
	}

	event.NewNodeCreatedPublisher(nats.Client()).Publish(event.NodeCreatedData{
		ProfileUrl:    node.ProfileUrl,
		LinkedSchemas: node.LinkedSchemas,
		Version:       *node.Version,
	})

	return &node, nil
}

func (s *nodesService) SetNodeValid(node node.Node) error {
	node.ID = cryptoutil.GetSHA256(node.ProfileUrl)
	node.Status = constant.Validated
	node.FailedReasons = &[]string{}
	if err := node.Update(); err != nil {
		return err
	}
	return nil
}

func (s *nodesService) SetNodeInValid(node node.Node) error {
	node.ID = cryptoutil.GetSHA256(node.ProfileUrl)
	node.Status = constant.ValidationFailed
	emptystr := ""
	node.ProfileHash = &emptystr
	node.LastValidated = dateutil.GetZeroValueUnix()
	if err := node.Update(); err != nil {
		return err
	}
	return nil
}

func (s *nodesService) SearchNode(query *node.NodeQuery) (node.Nodes, resterr.RestErr) {
	node := &node.Node{}
	result, err := node.Search(query)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *nodesService) DeleteNode(nodeId string) resterr.RestErr {
	node := &node.Node{ID: nodeId}
	if err := node.Delete(); err != nil {
		return err
	}
	return nil
}
