package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/elastic"
	mongodb "github.com/MurmurationsNetwork/MurmurationsServices/pkg/mongo"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/natsclient"
)

// PingHandler responds to ping requests by checking the status of MongoDB and NATS.
func PingHandler(c *gin.Context) {
	if !checkMongoDB(c) {
		return
	}

	if !checkNATS(c) {
		return
	}

	if !checkES(c) {
		return
	}

	c.String(http.StatusOK, "pong!")
}

// checkMongoDB performs a health check on MongoDB. Returns true if MongoDB is healthy.
func checkMongoDB(c *gin.Context) bool {
	client := mongodb.Client.GetClient()
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

// checkNATS checks if NATS is connected. Returns true if NATS is connected.
func checkNATS(c *gin.Context) bool {
	if natsclient.GetInstance() == nil {
		return true
	}

	if !natsclient.IsConnected() {
		c.String(http.StatusInternalServerError, "NATS is not connected")
		return false
	}

	return true
}

// checkES performs a health check on Elasticsearch. Returns true if Elasticsearch is healthy.
func checkES(c *gin.Context) bool {
	if elastic.Client.GetClient() == nil {
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
