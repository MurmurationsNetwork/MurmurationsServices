package usecase

import (
	"fmt"
	"net/http"

	"github.com/MurmurationsNetwork/MurmurationsServices/common/constant"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/cryptoutil"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/dateutil"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/event"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/httputil"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/nats"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/resterr"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/adapter/repository/db"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/entity/node"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/entity/query"
)

type NodeUsecase interface {
	AddNode(node *node.Node) (*node.Node, resterr.RestErr)
	GetNode(nodeID string) (*node.Node, resterr.RestErr)
	SetNodeValid(node *node.Node) error
	SetNodeInvalid(node *node.Node) error
	Search(query *query.EsQuery) (*query.QueryResults, resterr.RestErr)
	Delete(nodeID string) resterr.RestErr
}

type nodeUsecase struct {
	nodeRepo db.NodeRepository
}

func NewNodeService(nodeRepo db.NodeRepository) NodeUsecase {
	return &nodeUsecase{
		nodeRepo: nodeRepo,
	}
}

func (s *nodeUsecase) AddNode(node *node.Node) (*node.Node, resterr.RestErr) {
	if err := node.Validate(); err != nil {
		return nil, err
	}

	node.ID = cryptoutil.GetSHA256(node.ProfileURL)
	node.Status = constant.NodeStatus.Received
	node.CreatedAt = dateutil.GetNowUnix()

	if err := s.nodeRepo.Add(node); err != nil {
		return nil, err
	}

	logger.Info("Receiving a new node: " + node.ProfileURL)
	event.NewNodeCreatedPublisher(nats.Client.Client()).Publish(event.NodeCreatedData{
		ProfileURL: node.ProfileURL,
		Version:    *node.Version,
	})

	return node, nil
}

func (s *nodeUsecase) GetNode(nodeID string) (*node.Node, resterr.RestErr) {
	node := node.Node{ID: nodeID}
	err := s.nodeRepo.Get(&node)
	if err != nil {
		return nil, err
	}
	return &node, nil
}

func (s *nodeUsecase) SetNodeValid(node *node.Node) error {
	node.ID = cryptoutil.GetSHA256(node.ProfileURL)
	node.Status = constant.NodeStatus.Validated
	node.FailureReasons = &[]string{}

	if err := s.nodeRepo.Update(node); err != nil {
		return err
	}
	return nil
}

func (s *nodeUsecase) SetNodeInvalid(node *node.Node) error {
	node.ID = cryptoutil.GetSHA256(node.ProfileURL)
	node.Status = constant.NodeStatus.ValidationFailed
	emptystr := ""
	node.ProfileHash = &emptystr
	lastValidated := dateutil.GetZeroValueUnix()
	node.LastValidated = &lastValidated

	if err := s.nodeRepo.Update(node); err != nil {
		return err
	}
	return nil
}

func (s *nodeUsecase) Search(query *query.EsQuery) (*query.QueryResults, resterr.RestErr) {
	result, err := s.nodeRepo.Search(query)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *nodeUsecase) Delete(nodeID string) resterr.RestErr {
	node := &node.Node{ID: nodeID}

	if err := s.nodeRepo.Get(node); err != nil {
		return err
	}

	// TODO: Maybe we should avoid network requests in the index server?
	resp, err := httputil.Get(node.ProfileURL)
	if err != nil {
		return resterr.NewBadRequestError(fmt.Sprintf("Error when trying to get the information from node_id: %s", nodeID))
	}

	if resp.StatusCode == http.StatusOK {
		return resterr.NewBadRequestError(fmt.Sprintf("Profile still exists for node_id: %s", nodeID))
	}

	if resp.StatusCode == http.StatusNotFound {
		err := s.nodeRepo.Delete(node)
		if err != nil {
			return err
		}
		return nil
	}

	return resterr.NewBadRequestError(fmt.Sprintf("Node %s retrun status code %d", nodeID, resp.StatusCode))
}
