package queue

import (
	"log"

	"github.com/MurmurationsNetwork/MurmurationsServices/common/nats_utils"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/validation/internal/datasources/nats"
)

func Listen() {
	nodeCreatedListener := nats_utils.NewListener(nats.Client(), "node:created", "indexer-svc-qgroup")
	err := nodeCreatedListener.Listen()
	if err != nil {
		log.Fatal(err)
	}
}
