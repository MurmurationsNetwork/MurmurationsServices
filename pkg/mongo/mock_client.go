package mongo

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mockClient struct {
	client *mongo.Client
}

func (c *mockClient) FindOne(
	_ string,
	_ primitive.M,
) *mongo.SingleResult {
	return nil
}

func (c *mockClient) Count(
	_ string,
	_ primitive.M,
) (int64, error) {
	return 0, nil
}

func (c *mockClient) InsertOne(
	_ string,
	_ interface{},
	_ ...*options.InsertOneOptions,
) (*mongo.InsertOneResult, error) {
	return &mongo.InsertOneResult{}, nil
}

func (c *mockClient) FindOneAndUpdate(
	_ string,
	_ primitive.M,
	_ primitive.M,
	_ ...*options.FindOneAndUpdateOptions,
) (*mongo.SingleResult, error) {
	return &mongo.SingleResult{}, nil
}

func (c *mockClient) Find(
	_ string,
	_ primitive.M,
	_ ...*options.FindOptions,
) (*mongo.Cursor, error) {
	return &mongo.Cursor{}, nil
}

func (c *mockClient) DeleteOne(_ string, _ primitive.M) error {
	return nil
}

func (c *mockClient) DeleteMany(_ string, _ primitive.M) error {
	return nil
}

func (c *mockClient) Ping() error {
	return nil
}

func (c *mockClient) Disconnect() {
}

func (c *mockClient) GetClient() *mongo.Client {
	return c.client
}

func (c *mockClient) setClient(_ *mongo.Client, _ string) {
}
