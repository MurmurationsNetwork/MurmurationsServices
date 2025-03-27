package mongo

import (
	"context"
	"os"
	"time"

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

	CreateUniqueIndex(
		collection, indexName string,
		opts ...*options.CreateIndexesOptions,
	) error
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

	clientOptions := options.Client().ApplyURI(url)
	clientOptions.SetConnectTimeout(10 * time.Second)
	clientOptions.SetMaxConnIdleTime(5 * time.Minute)

	if os.Getenv("APP_ENV") != "test" {
		var err error
		client, err = mongo.Connect(context.Background(), clientOptions)
		if err != nil {
			return err
		}
	}

	Client.setClient(client, dbName)

	return nil
}
