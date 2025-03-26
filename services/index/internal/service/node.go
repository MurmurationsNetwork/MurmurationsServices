package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/constant"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/cryptoutil"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/dateutil"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/httputil"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/messaging"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/profile/profilehasher"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/config"
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

// SetNodeValid sets a node as valid, updates its status, and indexes it in the
// repositories.
func (s *nodeService) SetNodeValid(node *model.Node) error {
	// Prepare the node.
	node.ID = cryptoutil.ComputeSHA256(node.ProfileURL)
	node.SetStatusValidated()
	node.ResetFailureReasons()

	// Retrieve the old node from the database.
	oldNode, err := s.mongoRepo.GetByID(node.ID)
	if err != nil && !errors.As(err, &index.NotFoundError{}) {
		return err
	}

	if s.isProfileHashUnchanged(node, oldNode) {
		logger.Info(
			fmt.Sprintf(
				"Node with profile hash '%s' is unchanged.",
				*node.ProfileHash,
			),
		)
		// Reposted node with no changes has same last_updated timestamp.
		node.LastUpdated = oldNode.LastUpdated
	}

	// Update the node in MongoDB.
	if err := s.mongoRepo.Update(node); err != nil {
		return err
	}

	profile := model.NewProfile(node.ProfileStr)
	if err := profile.Update(node.ProfileURL, node.LastUpdated); err != nil {
		return err
	}

	profileJSON := profile.GetJSON()
	if node.Expires != nil {
		profileJSON["expires"] = *node.Expires
	}

	// Update Elastic Search.
	if err := s.elasticRepo.IndexByID(node.ID, profileJSON); err != nil {
		errMsg := fmt.Sprintf(
			"Error indexing node ID '%s' in Elastic repository.",
			node.ID,
		)
		logger.Error(errMsg, err)

		node.SetStatusPostFailed()
		if mongoErr := s.mongoRepo.Update(node); mongoErr != nil {
			logger.Error(
				"Failed to update node in MongoDB after Elastic indexing failure.",
				mongoErr,
			)
		}
		return err
	}

	// Set final status and update.
	node.SetStatusPosted()
	return s.mongoRepo.Update(node)
}

// isProfileHashUnchanged checks if the profile hash of the new node matches
// the old node. It returns true if the hashes are the same.
func (s *nodeService) isProfileHashUnchanged(
	newNode *model.Node,
	oldNode *model.Node,
) bool {
	if oldNode == nil {
		return false
	}

	newHash, _ := profilehasher.NewFromString(newNode.ProfileStr,
		config.Values.Library.InternalURL).Hash()

	return newHash != "" && oldNode.ProfileHash != nil &&
		*oldNode.ProfileHash == newHash
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
	if err := validateProfileURL(node.ProfileURL); err != nil {
		return nil, err
	}

	node.ID = cryptoutil.ComputeSHA256(node.ProfileURL)

	oldNode, err := s.mongoRepo.GetByID(node.ID)
	if err != nil && !errors.As(err, &index.NotFoundError{}) {
		return nil, err
	}

	// If oldNode is not nil and its status is 'Deleted', check if ProfileURL is valid.
	if oldNode != nil && oldNode.Status == constant.NodeStatus.Deleted {
		if !httputil.IsValidURL(node.ProfileURL) {
			return oldNode, nil
		}
	}

	node.Status = constant.NodeStatus.Received
	node.CreatedAt = dateutil.GetNowUnix()

	if err := s.mongoRepo.Add(node); err != nil {
		return nil, err
	}

	err = messaging.Publish(messaging.NodeCreated, messaging.NodeCreatedData{
		ProfileURL: node.ProfileURL,
		Version:    *node.Version,
	})
	if err != nil {
		return nil, err
	}

	return node, nil
}

func validateProfileURL(url string) error {
	if url == "" {
		return index.ValidationError{
			Field:  "ProfileURL",
			Reason: "The `profile_url` property is required.",
		}
	}
	if len(url) > 2000 {
		return index.ValidationError{
			Field:  "ProfileURL",
			Reason: "The `profile_url` property cannot exceed 2000 characters.",
		}
	}
	return nil
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

// Delete deletes a node based on its ID. It checks a feature toggle to decide
// whether to bypass the check for the profile URL's existence.
func (s *nodeService) Delete(nodeID string) (string, error) {
	node, err := s.mongoRepo.GetByID(nodeID)
	if err != nil {
		return "", err
	}

	if config.Values.FeatureToggles["SkipProfileURLCheckOnDelete"] {
		return s.proceedWithDeletion(node)
	}

	if err := s.checkProfileURL(node); err != nil {
		return "", err
	}

	return s.proceedWithDeletion(node)
}

// checkProfileURL checks the profile URL's existence, content type.
func (s *nodeService) checkProfileURL(node *model.Node) error {
	resp, err := httputil.Get(node.ProfileURL)
	if err != nil {
		// This error message better reflects that there was an issue with the
		// HTTP request, not that the URL is non-existent.
		return index.DeleteNodeError{
			Message: "HTTP Request Failed",
			Detail: fmt.Sprintf(
				"Error making HTTP request to %s: %s",
				node.ProfileURL,
				err,
			),
			ProfileURL: node.ProfileURL,
			ErrorCode:  index.ErrorHTTPRequestFailed,
		}
	}
	defer resp.Body.Close()

	// Check if the response is JSON.
	var bodyJSON interface{}
	isJSON := true
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&bodyJSON)
	if err != nil {
		isJSON = false
	}

	// If the profile URL doesn't return a 200 OK status, or if the response is
	// not JSON, then we consider the profile URL invalid, or non-existent for
	// our purposes and return nil, indicating that it's safe to proceed with
	// deletion.
	if resp.StatusCode != http.StatusOK || !isJSON {
		return nil
	}

	return index.DeleteNodeError{
		Message:    "Profile Still Exists",
		Detail:     fmt.Sprintf("Profile URL %s still exists", node.ProfileURL),
		ProfileURL: node.ProfileURL,
		ErrorCode:  index.ErrorProfileStillExists,
	}
}

// proceedWithDeletion contains the logic to delete the node from the database.
func (s *nodeService) proceedWithDeletion(node *model.Node) (string, error) {
	var err error

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
