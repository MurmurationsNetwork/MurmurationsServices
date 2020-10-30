package node

import (
	"context"
	"errors"

	"github.com/MurmurationsNetwork/MurmurationsServices/common/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/mongoutil"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/resterr"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/datasources/mongo/nodes_db"
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

	return nil
}

func (node *Node) Search(query *NodeQuery) (Nodes, resterr.RestErr) {
	filter := bson.M{
		"linkedSchemas": query.Schema,
		"lastChecked": bson.M{
			"$gte": query.LastChecked,
		},
	}

	cursor, err := nodes_db.Collection.Find(context.Background(), filter)
	if err != nil {
		logger.Error("error when trying to search nodes", err)
		return nil, resterr.NewInternalServerError("error when trying to search nodes", errors.New("database error"))
	}
	defer cursor.Close(context.Background())

	results := make([]Node, 0)
	for cursor.Next(context.Background()) {
		var node Node
		err := cursor.Decode(&node)
		if err != nil {
			logger.Error("error when trying to decode node indo a node struct", err)
			return nil, resterr.NewInternalServerError("error when trying to search nodes", errors.New("database error"))
		}
		results = append(results, node)
	}

	return results, nil
}

func (node *Node) Delete() resterr.RestErr {
	filter := bson.M{"_id": node.ID}

	_, err := nodes_db.Collection.DeleteOne(context.Background(), filter)
	if err != nil {
		logger.Error("error when trying to delete a node", err)
		return resterr.NewInternalServerError("error when trying to delete a node", errors.New("database error"))
	}

	return nil
}
