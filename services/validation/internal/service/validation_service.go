package service

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/MurmurationsNetwork/MurmurationsServices/common/cryptoutil"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/dateutil"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/event"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/httputil"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/validation/internal/datasource/nats"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/validation/internal/domain/node"
	"github.com/xeipuuv/gojsonschema"
)

var (
	ValidationService validationServiceInterface = &validationService{}
)

type validationServiceInterface interface {
	ValidateNode(node *node.Node)
}

type validationService struct{}

func (v *validationService) ValidateNode(node *node.Node) {
	document := gojsonschema.NewReferenceLoader(node.ProfileUrl)

	data, err := document.LoadJSON()
	linkedSchemas := data.(map[string]interface{})["linkedSchemas"].([]interface{})

	for _, linkedSchema := range linkedSchemas {
		// TODO: Wait for library.
		schemaURL := "https://raw.githubusercontent.com/MurmurationsNetwork/MurmurationsLibrary/master/schemas/" + linkedSchema.(string) + ".json"

		fmt.Println("==================================")
		fmt.Printf("schemaURL %+v \n", schemaURL)
		fmt.Println("==================================")

		schema, err := gojsonschema.NewSchema(gojsonschema.NewReferenceLoader(schemaURL))
		if err != nil {
			sendNodeValidationFailedEvent(node, []string{"Could not read from schema: " + schemaURL})
			return
		}

		result, err := schema.Validate(document)
		if err != nil {
			sendNodeValidationFailedEvent(node, []string{"Could not read from profile url: " + node.ProfileUrl})
			return
		}

		if !result.Valid() {
			failedReasons := parseResultError(result.Errors())
			sendNodeValidationFailedEvent(node, failedReasons)
			return
		}
	}

	jsonByte, err := httputil.GetByte(node.ProfileUrl)
	buffer := bytes.Buffer{}
	json.Compact(&buffer, jsonByte)
	if err != nil {
		sendNodeValidationFailedEvent(node, []string{"Could not read from profile url: " + node.ProfileUrl})
		return
	}

	event.NewNodeValidatedPublisher(nats.Client()).Publish(event.NodeValidatedData{
		ProfileUrl:  node.ProfileUrl,
		ProfileHash: cryptoutil.GetSHA256(buffer.String()),
		ProfileStr:  buffer.String(),
		LastChecked: dateutil.GetNowUnix(),
		Version:     node.Version,
	})
}

func parseResultError(resultErrors []gojsonschema.ResultError) []string {
	failedReasons := make([]string, 0)
	for _, desc := range resultErrors {
		failedReasons = append(failedReasons, desc.String())
	}
	return failedReasons
}

func sendNodeValidationFailedEvent(node *node.Node, failedReasons []string) {
	event.NewNodeValidationFailedPublisher(nats.Client()).Publish(event.NodeValidationFailedData{
		ProfileUrl:    node.ProfileUrl,
		FailedReasons: failedReasons,
		Version:       node.Version,
	})
}
