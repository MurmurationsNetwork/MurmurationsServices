package usecase

import (
	"encoding/json"
	"fmt"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/constant"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/cryptoutil"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/dateutil"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/event"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/httputil"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/jsonapi"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/nats"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/adapter/repository/db"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/entity"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/entity/query"
	"net/http"
)

type NodeUsecase interface {
	AddNode(node *entity.Node) (*entity.Node, []jsonapi.Error)
	GetNode(nodeID string) (*entity.Node, []jsonapi.Error)
	SetNodeValid(node *entity.Node) error
	SetNodeInvalid(node *entity.Node) error
	Search(query *query.EsQuery) (*query.QueryResults, []jsonapi.Error)
	Delete(nodeID string) (string, []jsonapi.Error)
	Export(query *query.EsBlockQuery) (*query.BlockQueryResults, []jsonapi.Error)
	GetNodes(query *query.EsQuery) (*query.MapQueryResults, []jsonapi.Error)
}

type nodeUsecase struct {
	nodeRepo db.NodeRepository
}

func NewNodeService(nodeRepo db.NodeRepository) NodeUsecase {
	return &nodeUsecase{
		nodeRepo: nodeRepo,
	}
}

func (s *nodeUsecase) AddNode(node *entity.Node) (*entity.Node, []jsonapi.Error) {
	if node.ProfileURL == "" {
		return nil, jsonapi.NewError([]string{"Missing Required Property"}, []string{"The `profile_url` property is required."}, nil, []int{http.StatusBadRequest})
	}

	node.ID = cryptoutil.GetSHA256(node.ProfileURL)

	// issue-425: if the node is already in the index and the status is deleted, we need to check the profile is valid or not
	oldNode, err := s.nodeRepo.GetNode(node.ID)
	if err != nil {
		return nil, err
	}
	if oldNode != nil && oldNode.Status == constant.NodeStatus.Deleted {
		isValid := httputil.IsValidURL(node.ProfileURL)
		if !isValid {
			return oldNode, nil
		}
	}

	node.Status = constant.NodeStatus.Received
	node.CreatedAt = dateutil.GetNowUnix()

	if err := s.nodeRepo.Add(node); err != nil {
		return nil, err
	}

	event.NewNodeCreatedPublisher(nats.Client.Client()).Publish(event.NodeCreatedData{
		ProfileURL: node.ProfileURL,
		Version:    *node.Version,
	})

	return node, nil
}

func (s *nodeUsecase) GetNode(nodeID string) (*entity.Node, []jsonapi.Error) {
	node, err := s.nodeRepo.Get(nodeID)
	if err != nil {
		return nil, err
	}
	return node, nil
}

func (s *nodeUsecase) SetNodeValid(node *entity.Node) error {
	node.ID = cryptoutil.GetSHA256(node.ProfileURL)
	node.Status = constant.NodeStatus.Validated
	node.FailureReasons = &[]jsonapi.Error{}

	if err := s.nodeRepo.Update(node); err != nil {
		return err
	}
	return nil
}

func (s *nodeUsecase) SetNodeInvalid(node *entity.Node) error {
	node.ID = cryptoutil.GetSHA256(node.ProfileURL)
	node.Status = constant.NodeStatus.ValidationFailed
	emptystr := ""
	node.ProfileHash = &emptystr
	lastUpdated := dateutil.GetZeroValueUnix()
	node.LastUpdated = &lastUpdated

	if err := s.nodeRepo.Update(node); err != nil {
		return err
	}
	return nil
}

func (s *nodeUsecase) Search(query *query.EsQuery) (*query.QueryResults, []jsonapi.Error) {
	result, err := s.nodeRepo.Search(query)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *nodeUsecase) Delete(nodeID string) (string, []jsonapi.Error) {
	node, getErr := s.nodeRepo.Get(nodeID)
	if getErr != nil {
		return "", getErr
	}

	// TODO: Maybe we should avoid network requests in the index server?
	resp, err := httputil.Get(node.ProfileURL)
	// defer here to avoid error
	defer resp.Body.Close()
	if err != nil {
		return node.ProfileURL, jsonapi.NewError([]string{"Profile URL Not Found"}, []string{fmt.Sprintf("There was an error when trying to reach %s to delete node_id: %s", node.ProfileURL, nodeID)}, nil, []int{http.StatusBadRequest})
	}

	// check the response is json or not (issue-266)
	var bodyJson interface{}
	isJson := true
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&bodyJson)
	if err != nil {
		isJson = false
	}

	if resp.StatusCode == http.StatusNotFound || isJson == false {
		if node.Status == constant.NodeStatus.Posted || node.Status == constant.NodeStatus.Deleted {
			err := s.nodeRepo.SoftDelete(node)
			if err != nil {
				return node.ProfileURL, err
			}
			return node.ProfileURL, nil
		} else {
			err := s.nodeRepo.Delete(node)
			if err != nil {
				return node.ProfileURL, err
			}
			return node.ProfileURL, nil
		}
	}

	if resp.StatusCode == http.StatusOK {
		return node.ProfileURL, jsonapi.NewError([]string{"Profile Still Exists"}, []string{fmt.Sprintf("The profile could not be deleted from the Index because it still exists at the profile_url: %s", node.ProfileURL)}, nil, []int{http.StatusBadRequest})
	}

	return node.ProfileURL, jsonapi.NewError([]string{"Node Status Code Error"}, []string{fmt.Sprintf("The node at %s returned the following status code: %d", node.ProfileURL, resp.StatusCode)}, nil, []int{http.StatusBadRequest})
}

func (s *nodeUsecase) Export(query *query.EsBlockQuery) (*query.BlockQueryResults, []jsonapi.Error) {
	result, err := s.nodeRepo.Export(query)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *nodeUsecase) GetNodes(query *query.EsQuery) (*query.MapQueryResults, []jsonapi.Error) {
	result, err := s.nodeRepo.GetNodes(query)
	if err != nil {
		return nil, err
	}
	return result, nil
}
