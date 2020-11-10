package node

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/MurmurationsNetwork/MurmurationsServices/common/constant"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/jsonutil"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/mongoutil"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/resterr"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/datasources/elasticsearch"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/datasources/mongo/nodes_db"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/domain/query"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ErrUpdate = errors.New("update error")
)

func (node *Node) Add() resterr.RestErr {
	filter := bson.M{"_id": node.ID}
	update := bson.M{"$set": node}
	opt := options.FindOneAndUpdate().SetUpsert(true)

	result, err := mongoutil.FindOneAndUpdate(nodes_db.Collection, filter, update, opt)
	if err != nil {
		logger.Error("error when trying to create a node", err)
		return resterr.NewInternalServerError("error when tying to add a node", errors.New("database error"))
	}

	var updated Node
	result.Decode(&updated)
	node.Version = updated.Version

	return nil
}

func (node *Node) Get() resterr.RestErr {
	filter := bson.M{"_id": node.ID}

	result := nodes_db.Collection.FindOne(context.Background(), filter)
	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return resterr.NewNotFoundError(fmt.Sprintf("Could not find node_id: %s", node.ID))
		}
		logger.Error("error when trying to find a node", result.Err())
		return resterr.NewInternalServerError("error when tying to find a node", errors.New("database error"))
	}

	err := result.Decode(&node)
	if err != nil {
		logger.Error("error when trying to parse database response", result.Err())
		return resterr.NewInternalServerError("error when tying to find a node", errors.New("database error"))
	}

	return nil
}

func (node *Node) Update() error {
	filter := bson.M{"_id": node.ID, "version": node.Version}
	// Unset the version to prevent setting it.
	node.Version = nil
	update := bson.M{"$set": node}

	_, err := mongoutil.FindOneAndUpdate(nodes_db.Collection, filter, update)
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
	if node.Status == constant.NodeStatus().Validated {
		profileJSON := jsonutil.ToJSON(node.ProfileStr)
		profileJSON["profile_url"] = node.ProfileURL
		profileJSON["last_validated"] = node.LastValidated

		_, err := elasticsearch.Client.IndexWithID(string(constant.ESIndex().Node), node.ID, profileJSON)
		if err != nil {
			// Fail to parse into ElasticSearch, set the statue to 'post_failed'.
			err = node.setPostFailed()
			if err != nil {
				return err
			}
		}

		// Successfully parse into ElasticSearch, set the statue to 'posted'.
		err = node.setPosted()
		if err != nil {
			return err
		}
	}

	if node.Status == constant.NodeStatus().ValidationFailed {
		err := elasticsearch.Client.Delete(string(constant.ESIndex().Node), node.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (node *Node) setPostFailed() error {
	node.Version = nil
	node.Status = constant.NodeStatus().PostFailed

	filter := bson.M{"_id": node.ID}
	update := bson.M{"$set": node}

	_, err := mongoutil.FindOneAndUpdate(nodes_db.Collection, filter, update)
	if err != nil {
		logger.Error("error when trying to update a node", err)
		return err
	}

	return nil
}

func (node *Node) setPosted() error {
	node.Version = nil
	node.Status = constant.NodeStatus().Posted

	filter := bson.M{"_id": node.ID}
	update := bson.M{"$set": node}

	_, err := mongoutil.FindOneAndUpdate(nodes_db.Collection, filter, update)
	if err != nil {
		logger.Error("error when trying to update a node", err)
		return err
	}

	return nil
}

func (node *Node) Search(q *query.EsQuery) (query.QueryResults, resterr.RestErr) {
	result, err := elasticsearch.Client.Search(string(constant.ESIndex().Node), q.Build())
	if err != nil {
		return nil, resterr.NewInternalServerError("error when trying to search documents", errors.New("database error"))
	}

	queryResults := make(query.QueryResults, result.TotalHits())
	for index, hit := range result.Hits.Hits {
		bytes, _ := hit.Source.MarshalJSON()
		var result query.QueryResult
		if err := json.Unmarshal(bytes, &result); err != nil {
			return nil, resterr.NewInternalServerError("error when trying to parse response", errors.New("database error"))
		}
		queryResults[index] = result
	}

	if len(queryResults) == 0 {
		return nil, resterr.NewNotFoundError("no items found matching given criteria")
	}

	return queryResults, nil
}

func (node *Node) Delete() resterr.RestErr {
	filter := bson.M{"_id": node.ID}

	// TODO: Abstract MongoDB operations.
	_, err := nodes_db.Collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return resterr.NewInternalServerError("error when trying to delete a node", errors.New("database error"))
	}
	err = elasticsearch.Client.Delete(string(constant.ESIndex().Node), node.ID)
	if err != nil {
		return resterr.NewInternalServerError("error when trying to delete a node", errors.New("database error"))
	}

	return nil
}
