package app

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/services/indexer/controllers/nodes"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/indexer/controllers/ping"
)

func mapUrls() {
	router.GET("/ping", ping.Ping)

	router.POST("/nodes", nodes.AddNode)
	router.GET("/nodes", nodes.GetNode)
	router.DELETE("/nodes/:node_id", nodes.DeleteNode)
}
