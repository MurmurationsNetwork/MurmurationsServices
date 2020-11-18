package nodes_db

import (
	"context"
	"os"
	"time"

	"github.com/MurmurationsNetwork/MurmurationsServices/common/logger"
	"github.com/cenkalti/backoff"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	Collection *mongo.Collection
	client     *mongo.Client
)

func init() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var err error
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err != nil {
		logger.Panic("error when trying to connect to MongoDB", err)
	}

	ping(client)

	Collection = client.Database("murmurations").Collection("nodes")
}

func Disconnect() {
	logger.Info("trying to disconnect from MongoDB")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Disconnect(ctx); err != nil {
		logger.Panic("error when trying to disconnect from MongoDB", err)
	}
}

func ping(client *mongo.Client) {
	op := func() error {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		err := client.Ping(ctx, readpref.Primary())
		if err != nil {
			return err
		}
		return nil
	}
	notify := func(err error, time time.Duration) {
		logger.Error("trying to re-connect MongoDB %s \n", err)
	}

	b := backoff.NewExponentialBackOff()
	b.MaxElapsedTime = 2 * time.Minute
	err := backoff.RetryNotify(op, b, notify)
	if err != nil {
		logger.Panic("error when trying to ping the MongoDB", err)
	}
}
