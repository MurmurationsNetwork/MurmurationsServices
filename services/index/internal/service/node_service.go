package service

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/common/constant"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/cryptoutil"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/dateutil"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/event"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/resterr"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/datasources/nats"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/domain/node"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/domain/query"
)

var (
	NodeService nodesServiceInterface = &nodesService{}
)

type nodesServiceInterface interface {
	AddNode(node node.Node) (*node.Node, resterr.RestErr)
	GetNode(nodeId string) (*node.Node, resterr.RestErr)
	SetNodeValid(node node.Node) error
	SetNodeInvalid(node node.Node) error
	Search(query *query.EsQuery) (query.QueryResults, resterr.RestErr)
	Delete(nodeId string) resterr.RestErr
}

type nodesService struct{}

func (s *nodesService) AddNode(node node.Node) (*node.Node, resterr.RestErr) {
	if err := node.Validate(); err != nil {
		return nil, err
	}

	node.ID = cryptoutil.GetSHA256(node.ProfileURL)
	node.Status = constant.NodeStatus().Received
	if err := node.Add(); err != nil {
		return nil, err
	}

	event.NewNodeCreatedPublisher(nats.Client()).Publish(event.NodeCreatedData{
		ProfileURL: node.ProfileURL,
		Version:    *node.Version,
	})

	return &node, nil
}

func (s *nodesService) GetNode(nodeId string) (*node.Node, resterr.RestErr) {
	node := node.Node{ID: nodeId}
	err := node.Get()
	if err != nil {
		return nil, err
	}
	return &node, nil
}

func (s *nodesService) SetNodeValid(node node.Node) error {
	node.ID = cryptoutil.GetSHA256(node.ProfileURL)
	node.Status = constant.NodeStatus().Validated
	node.FailureReasons = &[]string{}

	if err := node.Update(); err != nil {
		return err
	}
	return nil
}

func (s *nodesService) SetNodeInvalid(node node.Node) error {
	node.ID = cryptoutil.GetSHA256(node.ProfileURL)
	node.Status = constant.NodeStatus().ValidationFailed
	emptystr := ""
	node.ProfileHash = &emptystr
	lastValidated := dateutil.GetZeroValueUnix()
	node.LastValidated = &lastValidated

	if err := node.Update(); err != nil {
		return err
	}
	return nil
}

func (s *nodesService) Search(query *query.EsQuery) (query.QueryResults, resterr.RestErr) {
	dao := node.Node{}
	result, err := dao.Search(query)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *nodesService) Delete(nodeId string) resterr.RestErr {
	dao := node.Node{ID: nodeId}
	err := dao.Delete()
	if err != nil {
		return err
	}
	return nil
}
