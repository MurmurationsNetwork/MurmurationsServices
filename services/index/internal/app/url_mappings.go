package app

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/controller/http/node"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/controller/http/ping"
)

func mapUrls() {
	router.GET("/ping", ping.Ping)

	router.POST("/nodes", node.Add)
	router.GET("/nodes", node.Search)
	router.DELETE("/nodes/:node_id", node.Delete)
}
