package mongoutil

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func FindOneAndUpdate(c *mongo.Collection, filter primitive.M, update primitive.M, opts ...*options.FindOneAndUpdateOptions) (*mongo.SingleResult, error) {
	opts = append(opts, options.FindOneAndUpdate().SetReturnDocument(options.After))
	mergedOpt := options.MergeFindOneAndUpdateOptions(opts...)

	// Automatically increment the document version.
	update["$inc"] = bson.M{"version": 1}

	result := c.FindOneAndUpdate(context.Background(), filter, update, mergedOpt)
	if result.Err() != nil {
		return nil, result.Err()
	}

	return result, nil
}
