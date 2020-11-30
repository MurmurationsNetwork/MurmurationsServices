package mongo

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mockClient struct {
	client *mongo.Client
	db     *mongo.Database
}

func (c *mockClient) FindOne(collection string, filter primitive.M) *mongo.SingleResult {
	return nil
}

func (c *mockClient) FindOneAndUpdate(collection string, filter primitive.M, update primitive.M, opts ...*options.FindOneAndUpdateOptions) (*mongo.SingleResult, error) {
	return nil, nil
}

func (c *mockClient) DeleteOne(collection string, filter primitive.M) error {
	return nil
}

func (c *mockClient) Ping() error {
	return nil
}

func (c *mockClient) Disconnect() {
}

func (c *mockClient) setClient(client *mongo.Client, dbName string) {
}
