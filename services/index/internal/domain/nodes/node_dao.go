package nodes

import (
	"context"
	"errors"

	"github.com/MurmurationsNetwork/MurmurationsServices/common/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/rest_errors"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/datasources/mongo/nodes_db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ErrUpdate = errors.New("update error")
)

func (node *Node) Add() rest_errors.RestErr {
	filter := bson.M{"nodeId": node.NodeID}
	update := bson.M{"$set": node}
	opts := options.Update().SetUpsert(true)

	_, err := nodes_db.Collection.UpdateOne(context.Background(), filter, update, opts)
	if err != nil {
		logger.Error("error when trying to create a node", err)
		return rest_errors.NewInternalServerError("error when tying to add a node", errors.New("database error"))
	}

	return nil
}

func (node *Node) Update() error {
	filter := bson.M{"nodeId": node.NodeID}
	update := bson.M{"$set": node}
	opts := options.Update()

	_, err := nodes_db.Collection.UpdateOne(context.Background(), filter, update, opts)
	if err != nil {
		logger.Error("error when trying to update a node", err)
		return ErrUpdate
	}

	return nil
}

func (node *Node) Search(query *NodeQuery) (Nodes, rest_errors.RestErr) {
	filter := bson.M{
		"linkedSchemas": query.Schema,
		"lastValidated": bson.M{
			"$gte": query.LastValidated,
		},
	}

	cursor, err := nodes_db.Collection.Find(context.Background(), filter)
	if err != nil {
		logger.Error("error when trying to search nodes", err)
		return nil, rest_errors.NewInternalServerError("error when trying to search nodes", errors.New("database error"))
	}
	defer cursor.Close(context.Background())

	results := make([]Node, 0)
	for cursor.Next(context.Background()) {
		var node Node
		err := cursor.Decode(&node)
		if err != nil {
			logger.Error("error when trying to decode node indo a node struct", err)
			return nil, rest_errors.NewInternalServerError("error when trying to search nodes", errors.New("database error"))
		}
		results = append(results, node)
	}

	return results, nil
}

func (node *Node) Delete() rest_errors.RestErr {
	filter := bson.M{"nodeId": node.NodeID}

	_, err := nodes_db.Collection.DeleteOne(context.Background(), filter)
	if err != nil {
		logger.Error("error when trying to delete a node", err)
		return rest_errors.NewInternalServerError("error when trying to delete a node", errors.New("database error"))
	}

	return nil
}
