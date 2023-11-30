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
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/nats"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/index"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/model"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/repository/es"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/repository/mongo"
)

// NodeService is an interface that defines operations on nodes.
type NodeService interface {
	AddNode(node *model.Node) (*model.Node, error)
	GetNode(nodeID string) (*model.Node, error)
	SetNodeValid(node *model.Node) error
	SetNodeInvalid(node *model.Node) error
	Search(query *es.Query) (*es.QueryResults, error)
	Delete(nodeID string) (string, error)
	Export(query *es.BlockQuery) (*es.BlockQueryResults, error)
	GetNodes(query *es.Query) (*es.MapQueryResults, error)
}

type nodeService struct {
	mongoRepo   mongo.NodeRepository
	elasticRepo es.NodeRepository
}

// NewNodeService creates a new instance of NodeService.
func NewNodeService(
	mongoRepo mongo.NodeRepository,
	elasticRepo es.NodeRepository,

) NodeService {
	return &nodeService{
		mongoRepo:   mongoRepo,
		elasticRepo: elasticRepo,
	}
}

// SetNodeValid sets a node as valid.
func (s *nodeService) SetNodeValid(node *model.Node) error {
	node.ID = cryptoutil.ComputeSHA256(node.ProfileURL)
	node.SetStatusValidated()
	node.ResetFailureReasons()

	profile := model.NewProfile(node.ProfileStr)
	if err := profile.Update(node.ProfileURL, node.LastUpdated); err != nil {
		return err
	}

	if err := s.mongoRepo.Update(node); err != nil {
		return err
	}

	if err := s.elasticRepo.IndexByID(node.ID, profile.GetJSON()); err != nil {
		errMsg := fmt.Sprintf(
			"Error indexing node ID '%s' in Elastic repository",
			node.ID,
		)
		logger.Error(errMsg, err)
		node.SetStatusPostFailed()
		return s.mongoRepo.Update(node)
	}

	node.SetStatusPosted()
	return s.mongoRepo.Update(node)
}

// SetNodeInvalid sets a node as invali.
func (s *nodeService) SetNodeInvalid(node *model.Node) error {
	node.ID = cryptoutil.ComputeSHA256(node.ProfileURL)
	node.Status = constant.NodeStatus.ValidationFailed
	emptystr := ""
	node.ProfileHash = &emptystr
	lastUpdated := dateutil.GetZeroValueUnix()
	node.LastUpdated = &lastUpdated

	if err := s.mongoRepo.Update(node); err != nil {
		return err
	}

	return s.elasticRepo.DeleteByID(node.ID)
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

	node.ID = cryptoutil.ComputeSHA256(node.ProfileURL)

	oldNode, err := s.mongoRepo.GetByID(node.ID)

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

	node.Status = constant.NodeStatus.Received
	node.CreatedAt = dateutil.GetNowUnix()

	if err := s.mongoRepo.Add(node); err != nil {
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
	node, err := s.mongoRepo.GetByID(nodeID)
	if err != nil {
		return nil, err
	}
	return node, nil
}

// Search performs a search operation based on the provided query.
func (s *nodeService) Search(query *es.Query) (*es.QueryResults, error) {
	result, err := s.elasticRepo.Search(query)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Delete deletes a node based on its ID.
func (s *nodeService) Delete(nodeID string) (string, error) {
	node, err := s.mongoRepo.GetByID(nodeID)
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
			if err := s.mongoRepo.SoftDelete(node); err != nil {
				return node.ProfileURL, err
			}
			err = s.elasticRepo.SoftDelete(node)
			return node.ProfileURL, err
		}

		if err = s.mongoRepo.Delete(node); err != nil {
			return node.ProfileURL, err
		}
		err = s.elasticRepo.DeleteByID(node.ID)
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
	query *es.BlockQuery,
) (*es.BlockQueryResults, error) {
	result, err := s.elasticRepo.Export(query)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Export exports nodes based on the provided query.
func (s *nodeService) GetNodes(
	query *es.Query,
) (*es.MapQueryResults, error) {
	result, err := s.elasticRepo.GetNodes(query)
	if err != nil {
		return nil, err
	}
	return result, nil
}
