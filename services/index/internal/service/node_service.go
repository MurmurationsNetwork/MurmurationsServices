package service

import (
	"fmt"

	"github.com/MurmurationsNetwork/MurmurationsServices/common/constant"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/cryptoutil"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/dateutil"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/event"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/httputil"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/nats"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/resterr"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/domain/node"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/domain/query"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/repository/noderepo"
)

var (
	NodeService nodesServiceInterface = &nodesService{}
)

type nodesServiceInterface interface {
	AddNode(node *node.Node) (*node.Node, resterr.RestErr)
	GetNode(nodeID string) (*node.Node, resterr.RestErr)
	SetNodeValid(node *node.Node) error
	SetNodeInvalid(node *node.Node) error
	Search(query *query.EsQuery) (*query.QueryResults, resterr.RestErr)
	Delete(nodeID string) resterr.RestErr
}

type nodesService struct{}

func (s *nodesService) AddNode(node *node.Node) (*node.Node, resterr.RestErr) {
	if err := node.Validate(); err != nil {
		return nil, err
	}

	node.ID = cryptoutil.GetSHA256(node.ProfileURL)
	node.Status = constant.NodeStatus.Received

	if err := noderepo.Node.Add(node); err != nil {
		return nil, err
	}

	event.NewNodeCreatedPublisher(nats.Client.Client()).Publish(event.NodeCreatedData{
		ProfileURL: node.ProfileURL,
		Version:    *node.Version,
	})

	return node, nil
}

func (s *nodesService) GetNode(nodeID string) (*node.Node, resterr.RestErr) {
	node := node.Node{ID: nodeID}
	err := noderepo.Node.Get(&node)
	if err != nil {
		return nil, err
	}
	return &node, nil
}

func (s *nodesService) SetNodeValid(node *node.Node) error {
	node.ID = cryptoutil.GetSHA256(node.ProfileURL)
	node.Status = constant.NodeStatus.Validated
	node.FailureReasons = &[]string{}

	if err := noderepo.Node.Update(node); err != nil {
		return err
	}
	return nil
}

func (s *nodesService) SetNodeInvalid(node *node.Node) error {
	node.ID = cryptoutil.GetSHA256(node.ProfileURL)
	node.Status = constant.NodeStatus.ValidationFailed
	emptystr := ""
	node.ProfileHash = &emptystr
	lastValidated := dateutil.GetZeroValueUnix()
	node.LastValidated = &lastValidated

	if err := noderepo.Node.Update(node); err != nil {
		return err
	}
	return nil
}

func (s *nodesService) Search(query *query.EsQuery) (*query.QueryResults, resterr.RestErr) {
	result, err := noderepo.Node.Search(query)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *nodesService) Delete(nodeID string) resterr.RestErr {
	node := &node.Node{ID: nodeID}

	err := noderepo.Node.Get(node)
	if err != nil {
		return err
	}

	// TODO: Maybe we should avoid network requests in the index server?
	isValid := httputil.IsValidURL(node.ProfileURL)
	if isValid {
		return resterr.NewBadRequestError(fmt.Sprintf("Profile still exists for node_id: %s", nodeID))
	}

	err = noderepo.Node.Delete(node)
	if err != nil {
		return err
	}
	return nil
}
