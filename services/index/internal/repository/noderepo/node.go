package noderepo

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/MurmurationsNetwork/MurmurationsServices/common/constant"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/elastic"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/jsonutil"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/mongo"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/pagination"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/resterr"
	model "github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/domain/node"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/domain/query"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	Node nodeInterface
)

type nodeInterface interface {
	Add(node *model.Node) resterr.RestErr
	Get(node *model.Node) resterr.RestErr
	Update(node *model.Node) error
	Search(q *query.EsQuery) (*query.QueryResults, resterr.RestErr)
	Delete(node *model.Node) resterr.RestErr
}

type node struct{}

func init() {
	if os.Getenv("ENV") == "test" {
		Node = &mockNode{}
		return
	}
	Node = &node{}
}

func (dao *node) Add(node *model.Node) resterr.RestErr {
	filter := bson.M{"_id": node.ID}
	update := bson.M{"$set": node}
	opt := options.FindOneAndUpdate().SetUpsert(true)

	result, err := mongo.Client.FindOneAndUpdate(constant.MongoIndex.Node, filter, update, opt)
	if err != nil {
		logger.Error("Error when trying to create a node", err)
		return resterr.NewInternalServerError("Error when trying to add a node.", errors.New("database error"))
	}

	var updated model.Node
	result.Decode(&updated)
	node.Version = updated.Version

	return nil
}

func (dao *node) Get(node *model.Node) resterr.RestErr {
	filter := bson.M{"_id": node.ID}

	result := mongo.Client.FindOne(constant.MongoIndex.Node, filter)
	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return resterr.NewNotFoundError(fmt.Sprintf("Could not find node_id: %s", node.ID))
		}
		logger.Error("Error when trying to find a node", result.Err())
		return resterr.NewInternalServerError("Error when trying to find a node.", errors.New("database error"))
	}

	err := result.Decode(&node)
	if err != nil {
		logger.Error("Error when trying to parse database response", result.Err())
		return resterr.NewInternalServerError("Error when trying to find a node.", errors.New("database error"))
	}

	return nil
}

func (dao *node) Update(node *model.Node) error {
	filter := bson.M{"_id": node.ID, "version": node.Version}
	// Unset the version to prevent setting it.
	node.Version = nil
	update := bson.M{"$set": node}

	_, err := mongo.Client.FindOneAndUpdate(constant.MongoIndex.Node, filter, update)
	if err != nil {
		// Update the document only if the version matches.
		// If the version does not match, it's an expected concurrent issue.
		if err == mongo.ErrNoDocuments {
			return nil
		}
		logger.Error("error when trying to update a node", err)
		return ErrUpdate
	}

	// NOTE: Maybe it's better to conver into another event?
	if node.Status == constant.NodeStatus.Validated {
		profileJSON := jsonutil.ToJSON(node.ProfileStr)
		profileJSON["profile_url"] = node.ProfileURL
		profileJSON["last_validated"] = node.LastValidated

		_, err := elastic.Client.IndexWithID(constant.ESIndex.Node, node.ID, profileJSON)
		if err != nil {
			// Fail to parse into ElasticSearch, set the statue to 'post_failed'.
			err = dao.setPostFailed(node)
			if err != nil {
				return err
			}
		} else {
			// Successfully parse into ElasticSearch, set the statue to 'posted'.
			err = dao.setPosted(node)
			if err != nil {
				return err
			}
		}
	}

	if node.Status == constant.NodeStatus.ValidationFailed {
		err := elastic.Client.Delete(constant.ESIndex.Node, node.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (dao *node) setPostFailed(node *model.Node) error {
	node.Version = nil
	node.Status = constant.NodeStatus.PostFailed

	filter := bson.M{"_id": node.ID}
	update := bson.M{"$set": node}

	_, err := mongo.Client.FindOneAndUpdate(constant.MongoIndex.Node, filter, update)
	if err != nil {
		logger.Error("error when trying to update a node", err)
		return err
	}

	return nil
}

func (dao *node) setPosted(node *model.Node) error {
	node.Version = nil
	node.Status = constant.NodeStatus.Posted

	filter := bson.M{"_id": node.ID}
	update := bson.M{"$set": node}

	_, err := mongo.Client.FindOneAndUpdate(constant.MongoIndex.Node, filter, update)
	if err != nil {
		logger.Error("error when trying to update a node", err)
		return err
	}

	return nil
}

func (dao *node) Search(q *query.EsQuery) (*query.QueryResults, resterr.RestErr) {
	result, err := elastic.Client.Search(constant.ESIndex.Node, q.Build())
	if err != nil {
		return nil, resterr.NewInternalServerError("Error when trying to search documents.", errors.New("database error"))
	}

	queryResults := make([]query.QueryResult, 0)
	for _, hit := range result.Hits.Hits {
		bytes, _ := hit.Source.MarshalJSON()
		var result query.QueryResult
		if err := json.Unmarshal(bytes, &result); err != nil {
			return nil, resterr.NewInternalServerError("Error when trying to parse response.", errors.New("database error"))
		}
		queryResults = append(queryResults, result)
	}

	if len(queryResults) == 0 {
		return nil, resterr.NewNotFoundError("No items found matching given criteria.")
	}

	return &query.QueryResults{
		Result:          queryResults,
		NumberOfResults: result.Hits.TotalHits,
		TotalPages:      pagination.TotalPages(result.Hits.TotalHits, q.PageSize),
	}, nil
}

func (dao *node) Delete(node *model.Node) resterr.RestErr {
	filter := bson.M{"_id": node.ID}

	err := mongo.Client.DeleteOne(constant.MongoIndex.Node, filter)
	if err != nil {
		return resterr.NewInternalServerError("Error when trying to delete a node.", errors.New("database error"))
	}
	err = elastic.Client.Delete(constant.ESIndex.Node, node.ID)
	if err != nil {
		return resterr.NewInternalServerError("Error when trying to delete a node.", errors.New("database error"))
	}

	return nil
}
