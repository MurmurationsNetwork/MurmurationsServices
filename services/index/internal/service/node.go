package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/constant"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/cryptoutil"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/dateutil"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/event"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/httputil"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/jsonapi"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/nats"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/index"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/model"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/model/query"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/repository/db"
)

// NodeService is an interface that defines operations on nodes.
type NodeService interface {
	AddNode(node *model.Node) (*model.Node, error)
	GetNode(nodeID string) (*model.Node, error)
	SetNodeValid(node *model.Node) error
	SetNodeInvalid(node *model.Node) error
	Search(query *query.EsQuery) (*query.Results, error)
	Delete(nodeID string) (string, error)
	Export(
		query *query.EsBlockQuery,
	) (*query.BlockQueryResults, error)
	GetNodes(query *query.EsQuery) (*query.MapQueryResults, error)
}

type nodeService struct {
	nodeRepo db.NodeRepository
}

// NewNodeService creates a new instance of NodeService.
func NewNodeService(nodeRepo db.NodeRepository) NodeService {
	return &nodeService{
		nodeRepo: nodeRepo,
	}
}

// SetNodeValid sets a node as valid.
func (s *nodeService) SetNodeValid(node *model.Node) error {
	node.ID = cryptoutil.GetSHA256(node.ProfileURL)
	node.Status = constant.NodeStatus.Validated
	node.FailureReasons = &[]jsonapi.Error{}
	return s.nodeRepo.Update(node)
}

// SetNodeInvalid sets a node as invali.
func (s *nodeService) SetNodeInvalid(node *model.Node) error {
	node.ID = cryptoutil.GetSHA256(node.ProfileURL)
	node.Status = constant.NodeStatus.ValidationFailed
	emptystr := ""
	node.ProfileHash = &emptystr
	lastUpdated := dateutil.GetZeroValueUnix()
	node.LastUpdated = &lastUpdated

	return s.nodeRepo.Update(node)
}

// AddNode adds a new node to the system.
func (s *nodeService) AddNode(
	node *model.Node,
) (*model.Node, error) {
	if node.ProfileURL == "" {
		return nil, index.ValidationError{
			Field:  "ProfileURL",
			Reason: "The `profile_url` property is required.",
		}
	}

	node.ID = cryptoutil.GetSHA256(node.ProfileURL)

	oldNode, err := s.nodeRepo.GetNode(node.ID)

	// If the error is not a NotFoundError, return the error.
	if err != nil && !errors.As(err, &index.NotFoundError{}) {
		return nil, err
	}
	// Handle the case where oldNode is found and its status is Deleted.
	if err == nil && oldNode.Status == constant.NodeStatus.Deleted {
		isValid := httputil.IsValidURL(node.ProfileURL)
		if !isValid {
			return oldNode, nil
		}
	}
	// Handle the case where oldNode is found and profile hash is the same.
	if err == nil {
		jsonStr, err := httputil.GetJSONStr(node.ProfileURL)
		if err == nil {
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

// GetNode retrieves a node based on its ID.
func (s *nodeService) GetNode(nodeID string) (*model.Node, error) {
	node, err := s.nodeRepo.Get(nodeID)
	if err != nil {
		return nil, err
	}
	return node, nil
}

// Search performs a search operation based on the provided query.
func (s *nodeService) Search(query *query.EsQuery) (*query.Results, error) {
	result, err := s.nodeRepo.Search(query)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Delete deletes a node based on its ID.
func (s *nodeService) Delete(nodeID string) (string, error) {
	node, err := s.nodeRepo.Get(nodeID)
	if err != nil {
		return "", err
	}

	resp, err := httputil.Get(node.ProfileURL)
	if err != nil {
		return "", index.DeleteNodeError{
			Message: "Profile URL Not Found",
			Detail: fmt.Sprintf(
				"There was an error when trying to reach %s to delete node_id: %s",
				node.ProfileURL,
				nodeID,
			),
			ProfileURL: node.ProfileURL,
			NodeID:     node.ID,
		}
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

	// check profile_url has redirect or not (issue-516)
	hasRedirect, err := httputil.CheckRedirect(node.ProfileURL)
	if err != nil {
		return "", index.DeleteNodeError{
			Message: "Profile URL Cannot Be Checked",
			Detail: fmt.Sprintf(
				"There was an error when trying to reach %s to delete node_id: %s",
				node.ProfileURL,
				nodeID,
			),
			NodeID: node.ID,
		}
	}

	if resp.StatusCode == http.StatusNotFound || !isJSON || hasRedirect {
		if node.Status == constant.NodeStatus.Posted ||
			node.Status == constant.NodeStatus.Deleted {
			err := s.nodeRepo.SoftDelete(node)
			return node.ProfileURL, err
		}
		err := s.nodeRepo.Delete(node)
		return node.ProfileURL, err
	}

	if resp.StatusCode == http.StatusOK {
		return node.ProfileURL, index.DeleteNodeError{
			Message: "Profile Still Exists",
			Detail: fmt.Sprintf(
				"The profile could not be deleted from the Index because "+
					"it still exists at the profile_url: %s",
				node.ProfileURL,
			),
			ProfileURL: node.ProfileURL,
			NodeID:     node.ID,
		}
	}

	return node.ProfileURL, index.DeleteNodeError{
		Message: "Profile Still Exists",
		Detail: fmt.Sprintf(
			"The node at %s returned the following status code: %d",
			node.ProfileURL,
			resp.StatusCode,
		),
		ProfileURL: node.ProfileURL,
		NodeID:     node.ID,
	}
}

// Export exports nodes based on the provided query.
func (s *nodeService) Export(
	query *query.EsBlockQuery,
) (*query.BlockQueryResults, error) {
	result, err := s.nodeRepo.Export(query)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Export exports nodes based on the provided query.
func (s *nodeService) GetNodes(
	query *query.EsQuery,
) (*query.MapQueryResults, error) {
	result, err := s.nodeRepo.GetNodes(query)
	if err != nil {
		return nil, err
	}
	return result, nil
}
