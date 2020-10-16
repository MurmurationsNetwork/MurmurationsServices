package app

import "github.com/MurmurationsNetwork/MurmurationsServices/indexer/datasources/mongo/nodes_db"

func cleanup() {
	nodes_db.Disconnect()
}
