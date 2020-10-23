package events

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/common/events"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/validation/internal/datasources/nats"
)

const qgroup = "validation-svc-qgroup"

var HandleNodeCreated = events.NewNodeCreatedListener(nats.Client(), qgroup, events.DefaultMsgHandler())
