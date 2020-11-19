package app

import "github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/adapter/mongo/nodes_db"

func cleanup() {
	nodes_db.Disconnect()
}
