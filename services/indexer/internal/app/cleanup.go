package app

import "github.com/MurmurationsNetwork/MurmurationsServices/services/indexer/internal/datasources/mongo/nodes_db"

func cleanup() {
	nodes_db.Disconnect()
}
