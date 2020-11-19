package nats

import (
	"os"

	"github.com/MurmurationsNetwork/MurmurationsServices/common/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/natswrapper"
	"github.com/nats-io/stan.go"
)

var (
	client stan.Conn
)

func init() {
	var err error
	client, err = natswrapper.Connect(os.Getenv("NATS_CLUSTER_ID"), os.Getenv("NATS_CLIENT_ID"), os.Getenv("NATS_URL"))
	if err != nil {
		logger.Panic("error when trying to connect nats", err)
	}
}

func Client() stan.Conn {
	return client
}

func Disconnect() {
	logger.Info("trying to disconnect from NATS")
	client.Close()
}
