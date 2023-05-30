package usecase

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/constant"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/cryptoutil"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/dateutil"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/event"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/httputil"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/jsonapi"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/nats"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/adapter/repository/db"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/entity"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/entity/query"
)

type NodeUsecase interface {
	AddNode(node *entity.Node) (*entity.Node, []jsonapi.Error)
	GetNode(nodeID string) (*entity.Node, []jsonapi.Error)
	SetNodeValid(node *entity.Node) error
	SetNodeInvalid(node *entity.Node) error
	Search(query *query.EsQuery) (*query.Results, []jsonapi.Error)
	Delete(nodeID string) (string, []jsonapi.Error)
	Export(
		query *query.EsBlockQuery,
	) (*query.BlockQueryResults, []jsonapi.Error)
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

func (s *nodeUsecase) AddNode(
	node *entity.Node,
) (*entity.Node, []jsonapi.Error) {
	if node.ProfileURL == "" {
		return nil, jsonapi.NewError(
			[]string{"Missing Required Property"},
			[]string{"The `profile_url` property is required."},
			nil,
			[]int{http.StatusBadRequest},
		)
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

	// issue-451: if old profile hash is the same as new profile hash, we can directly return the old node
	if oldNode != nil {
		jsonStr, jsonErr := httputil.GetJSONStr(node.ProfileURL)
		// when the error is nil, it means we can get data from the profile url, then we need to check the hash
		if jsonErr == nil {
			newHash := cryptoutil.GetSHA256(jsonStr)
			if oldNode.ProfileHash != nil && *oldNode.ProfileHash == newHash {
				return oldNode, nil
			}
		}
	}

	node.Status = constant.NodeStatus.Received
	node.CreatedAt = dateutil.GetNowUnix()

	if err := s.nodeRepo.Add(node); err != nil {
		return nil, err
	}

	event.NewNodeCreatedPublisher(nats.Client.Client()).
		Publish(event.NodeCreatedData{
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

	return s.nodeRepo.Update(node)
}

func (s *nodeUsecase) SetNodeInvalid(node *entity.Node) error {
	node.ID = cryptoutil.GetSHA256(node.ProfileURL)
	node.Status = constant.NodeStatus.ValidationFailed
	emptystr := ""
	node.ProfileHash = &emptystr
	lastUpdated := dateutil.GetZeroValueUnix()
	node.LastUpdated = &lastUpdated

	return s.nodeRepo.Update(node)
}

func (s *nodeUsecase) Search(
	query *query.EsQuery,
) (*query.Results, []jsonapi.Error) {
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
	if err != nil {
		return node.ProfileURL, jsonapi.NewError(
			[]string{"Profile URL Not Found"},
			[]string{
				fmt.Sprintf(
					"There was an error when trying to reach %s to delete node_id: %s",
					node.ProfileURL,
					nodeID,
				),
			},
			nil,
			[]int{http.StatusBadRequest},
		)
	}
	defer resp.Body.Close()

	// check the response is json or not (issue-266)
	var bodyJSON interface{}
	isJSON := true
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&bodyJSON)
	if err != nil {
		isJSON = false
	}

	if resp.StatusCode == http.StatusNotFound || !isJSON {
		if node.Status == constant.NodeStatus.Posted ||
			node.Status == constant.NodeStatus.Deleted {
			err := s.nodeRepo.SoftDelete(node)
			return node.ProfileURL, err
		}
		err := s.nodeRepo.Delete(node)
		return node.ProfileURL, err
	}

	if resp.StatusCode == http.StatusOK {
		return node.ProfileURL, jsonapi.NewError(
			[]string{"Profile Still Exists"},
			[]string{
				fmt.Sprintf(
					"The profile could not be deleted from the Index because it still exists at the profile_url: %s",
					node.ProfileURL,
				),
			},
			nil,
			[]int{http.StatusBadRequest},
		)
	}

	return node.ProfileURL, jsonapi.NewError(
		[]string{"Node Status Code Error"},
		[]string{
			fmt.Sprintf(
				"The node at %s returned the following status code: %d",
				node.ProfileURL,
				resp.StatusCode,
			),
		},
		nil,
		[]int{http.StatusBadRequest},
	)
}

func (s *nodeUsecase) Export(
	query *query.EsBlockQuery,
) (*query.BlockQueryResults, []jsonapi.Error) {
	result, err := s.nodeRepo.Export(query)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *nodeUsecase) GetNodes(
	query *query.EsQuery,
) (*query.MapQueryResults, []jsonapi.Error) {
	result, err := s.nodeRepo.GetNodes(query)
	if err != nil {
		return nil, err
	}
	return result, nil
}
