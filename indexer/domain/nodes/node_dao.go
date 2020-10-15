package nodes

import (
	"context"

	"github.com/MurmurationsNetwork/MurmurationsServices/indexer/datasources/mongo/nodes_db"
	"github.com/MurmurationsNetwork/MurmurationsServices/utils/date_utils"
	"github.com/MurmurationsNetwork/MurmurationsServices/utils/hash_utils"
	"github.com/MurmurationsNetwork/MurmurationsServices/utils/mongo_utils"
	"github.com/MurmurationsNetwork/MurmurationsServices/utils/rest_errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (node *Node) Add() rest_errors.RestErr {
	node.NodeID = hash_utils.SHA256(node.ProfileUrl)
	node.LastValidated = date_utils.GetNowUnix()

	filter := bson.M{"nodeId": node.NodeID}
	update := bson.M{"$set": node}
	opts := options.Update().SetUpsert(true)

	_, err := nodes_db.Collection.UpdateOne(context.Background(), filter, update, opts)
	if err != nil {
		return mongo_utils.ParseError(err)
	}

	return nil
}

func (node *Node) Get() rest_errors.RestErr {
	return nil
}

func (node *Node) Search() rest_errors.RestErr {
	return nil
}

func (node *Node) Delete() rest_errors.RestErr {
	return nil
}
