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

const (
	indexNodes = "nodes"
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

	if node.Status == constant.Validated {
		profileJson := jsonutil.ToJSON(node.ProfileStr)
		profileJson["lastChecked"] = node.LastChecked

		fmt.Println("==================================")
		fmt.Printf("profileJson %+v \n", profileJson)
		fmt.Println("==================================")

		result, err := elasticsearch.Client.IndexWithID(indexNodes, node.ID, profileJson)
		if err != nil {
			return err
		}

		fmt.Println("==================================")
		fmt.Printf("result %+v \n", result)
		fmt.Println("==================================")
	}

	return nil
}

func (node *Node) Search(q *query.EsQuery) (query.QueryResults, resterr.RestErr) {
	result, err := elasticsearch.Client.Search(indexNodes, q.Build())
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
