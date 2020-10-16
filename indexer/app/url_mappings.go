package app

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/indexer/controllers/nodes"
	"github.com/MurmurationsNetwork/MurmurationsServices/indexer/controllers/ping"
)

func mapUrls() {
	router.GET("/ping", ping.Ping)

	router.POST("/nodes", nodes.Add)
	router.GET("/nodes/:node_id", nodes.Get)
	router.GET("/nodes", nodes.Search)
	router.DELETE("/nodes/:node_id", nodes.Delete)
}
