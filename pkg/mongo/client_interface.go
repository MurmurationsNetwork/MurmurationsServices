package mongo

import (
	"context"
	"os"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	Client mongoClientInterface
)

type mongoClientInterface interface {
	FindOne(collection string, filter primitive.M) *mongo.SingleResult
	Count(collection string, filter primitive.M) (int64, error)
	InsertOne(collection string, document interface{},
		opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
	FindOneAndUpdate(
		collection string,
		filter primitive.M,
		update primitive.M,
		opts ...*options.FindOneAndUpdateOptions,
	) (*mongo.SingleResult, error)
	Find(
		collection string,
		filter primitive.M,
		opts ...*options.FindOptions,
	) (*mongo.Cursor, error)
	DeleteOne(collection string, filter primitive.M) error
	DeleteMany(collection string, filter primitive.M) error

	Ping() error
	Disconnect()

	GetClient() *mongo.Client
	setClient(*mongo.Client, string)
}

func init() {
	if os.Getenv("APP_ENV") == "test" {
		Client = &mockClient{}
		return
	}
	Client = &mongoClient{}
}

func NewClient(url string, dbName string) error {
	var client *mongo.Client

	if os.Getenv("APP_ENV") != "test" {
		var err error
		client, err = mongo.Connect(
			context.Background(),
			options.Client().ApplyURI(url),
		)
		if err != nil {
			return err
		}
	}

	Client.setClient(client, dbName)

	return nil
}
