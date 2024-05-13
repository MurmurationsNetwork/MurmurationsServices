package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/elastic"
	mongodb "github.com/MurmurationsNetwork/MurmurationsServices/pkg/mongo"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/natsclient"
)

// PingHandler handles ping requests and checks the status of MongoDB, NATS, and Elasticsearch.
func PingHandler(c *gin.Context) {
	if !checkMongoDB(c) || !checkNATS(c) || !checkES(c) {
		return
	}
	c.String(http.StatusOK, "pong!")
}

// checkMongoDB verifies MongoDB's health. Returns true if healthy or not configured.
func checkMongoDB(c *gin.Context) bool {
	client := mongodb.Client.GetClient()
	// If there is no client, assume MongoDB is not configured and return true.
	if client == nil {
		return true
	}

	err := mongodb.Client.Ping()
	if err != nil {
		c.String(
			http.StatusInternalServerError,
			"error pinging MongoDB: %s",
			err,
		)
		return false
	}

	return true
}

// checkNATS verifies if NATS is connected. Returns true if connected or not configured.
func checkNATS(c *gin.Context) bool {
	if natsclient.GetInstance() == nil {
		// If there is no instance, assume NATS is not configured and return true.
		return true
	}

	if !natsclient.IsConnected() {
		c.String(http.StatusInternalServerError, "NATS is not connected")
		return false
	}

	return true
}

// checkES verifies Elasticsearch's health. Returns true if healthy or not configured.
func checkES(c *gin.Context) bool {
	if elastic.Client.GetClient() == nil {
		// If there is no client, assume Elasticsearch is not configured and return true.
		return true
	}

	err := elastic.Client.Ping()
	if err != nil {
		c.String(
			http.StatusInternalServerError,
			"Error pinging Elasticsearch: %s",
			err,
		)
		return false
	}

	return true
}
