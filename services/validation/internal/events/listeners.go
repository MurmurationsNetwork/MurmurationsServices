package events

import (
	"fmt"

	"github.com/MurmurationsNetwork/MurmurationsServices/common/events"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/validation/internal/datasources/nats"
	"github.com/nats-io/stan.go"
	"github.com/xeipuuv/gojsonschema"
)

const qgroup = "validation-svc-qgroup"

var HandleNodeCreated = events.NewNodeCreatedListener(nats.Client(), qgroup, func(msg *stan.Msg) {
	var nodeCreatedData events.NodeCreatedData

	document := gojsonschema.NewReferenceLoader(nodeCreatedData.ProfileUrl)

	for _, schemaURL := range nodeCreatedData.LinkedSchemas {
		schemaLoader := gojsonschema.NewReferenceLoader(schemaURL)
		result, err := gojsonschema.Validate(schemaLoader, document)
		if err != nil {
			// Internet error retry
			panic(err.Error())
		}
		if !result.Valid() {
			fmt.Printf("The document is not valid. see errors :\n")
			for _, desc := range result.Errors() {
				fmt.Printf("- %s\n", desc)
			}
			// Invalid document.
		}
	}

	// Valid document.
	msg.Ack()
})
