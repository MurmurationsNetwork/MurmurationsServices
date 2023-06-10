package mongo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/retry"
)

type mongoClient struct {
	client *mongo.Client
	db     *mongo.Database
}

func (c *mongoClient) FindOne(
	collection string,
	filter primitive.M,
) *mongo.SingleResult {
	return c.db.Collection(collection).FindOne(context.Background(), filter)
}

func (c *mongoClient) Count(
	collection string,
	filter primitive.M,
) (int64, error) {
	return c.db.Collection(collection).
		CountDocuments(context.Background(), filter)
}

func (c *mongoClient) InsertOne(collection string, document interface{},
	opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	mergedOpt := options.MergeInsertOneOptions(opts...)

	result, err := c.db.Collection(collection).
		InsertOne(context.Background(), document, mergedOpt)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *mongoClient) FindOneAndUpdate(
	collection string,
	filter primitive.M,
	update primitive.M,
	opts ...*options.FindOneAndUpdateOptions,
) (*mongo.SingleResult, error) {
	opts = append(
		opts,
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)
	mergedOpt := options.MergeFindOneAndUpdateOptions(opts...)

	// Automatically increment the document version.
	update["$inc"] = bson.M{"__v": 1}

	result := c.db.Collection(collection).
		FindOneAndUpdate(context.Background(), filter, update, mergedOpt)
	if result.Err() != nil {
		return nil, result.Err()
	}

	return result, nil
}

func (c *mongoClient) Find(
	collection string,
	filter primitive.M,
	opts ...*options.FindOptions,
) (*mongo.Cursor, error) {
	cur, err := c.db.Collection(collection).
		Find(context.Background(), filter, opts...)
	if err != nil {
		return nil, err
	}
	return cur, nil
}

func (c *mongoClient) DeleteOne(collection string, filter primitive.M) error {
	_, err := c.db.Collection(collection).
		DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}
	return nil
}

func (c *mongoClient) DeleteMany(collection string, filter primitive.M) error {
	_, err := c.db.Collection(collection).
		DeleteMany(context.Background(), filter)
	if err != nil {
		return err
	}
	return nil
}

func (c *mongoClient) Ping() error {
	operation := func() error {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		err := c.client.Ping(ctx, readpref.Primary())
		if err != nil {
			return err
		}
		return nil
	}
	err := retry.Do(operation)
	if err != nil {
		return err
	}

	return nil
}

func (c *mongoClient) Disconnect() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := c.client.Disconnect(ctx)
	if err != nil {
		logger.Error("Error when trying to disconnect from MongoDB", err)
	}
}

func (c *mongoClient) GetClient() *mongo.Client {
	return c.client
}

func (c *mongoClient) setClient(client *mongo.Client, dbName string) {
	c.client = client
	c.db = client.Database(dbName)
}
