package service

import (
	"bytes"
	"encoding/json"

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
	if err != nil {
		sendNodeValidationFailedEvent(node, []string{"Could not read from profile url: " + node.ProfileUrl})
		return
	}

	linkedSchemas, ok := getLinkedSchemas(data)
	if !ok {
		sendNodeValidationFailedEvent(node, []string{"Could not read linkedSchemas from profile url: " + node.ProfileUrl})
		return
	}

	failedReasons := validateAgainstSchemas(linkedSchemas, document)
	if len(failedReasons) != 0 {
		sendNodeValidationFailedEvent(node, failedReasons)
		return
	}

	jsonStr, err := getJSONStr(node.ProfileUrl)
	if err != nil {
		sendNodeValidationFailedEvent(node, []string{"Could not get JSON string from profile url: " + node.ProfileUrl})
		return
	}

	event.NewNodeValidatedPublisher(nats.Client()).Publish(event.NodeValidatedData{
		ProfileUrl:  node.ProfileUrl,
		ProfileHash: cryptoutil.GetSHA256(jsonStr),
		ProfileStr:  jsonStr,
		LastChecked: dateutil.GetNowUnix(),
		Version:     node.Version,
	})
}

func getLinkedSchemas(data interface{}) ([]string, bool) {
	json, ok := data.(map[string]interface{})
	if !ok {
		return nil, false
	}
	_, ok = json["linkedSchemas"]
	if !ok {
		return nil, false
	}
	arrInterface, ok := json["linkedSchemas"].([]interface{})
	if !ok {
		return nil, false
	}

	var linkedSchemas = make([]string, 0)

	for _, data := range arrInterface {
		linkedSchema, ok := data.(string)
		if !ok {
			return nil, false
		}
		linkedSchemas = append(linkedSchemas, linkedSchema)
	}

	return linkedSchemas, true
}

func validateAgainstSchemas(linkedSchemas []string, document gojsonschema.JSONLoader) []string {
	failedReasons := []string{}

	for _, linkedSchema := range linkedSchemas {
		// TODO: Wait for library.
		schemaURL := "https://raw.githubusercontent.com/MurmurationsNetwork/MurmurationsLibrary/master/schemas/" + linkedSchema + ".json"

		schema, err := gojsonschema.NewSchema(gojsonschema.NewReferenceLoader(schemaURL))
		if err != nil {
			failedReasons = append(failedReasons, "Could not read from schema: "+schemaURL)
			continue
		}

		result, err := schema.Validate(document)
		if err != nil {
			failedReasons = append(failedReasons, "error when trying to validate document: ", err.Error())
			continue
		}

		if !result.Valid() {
			failedReasons = append(failedReasons, parseValidateError(linkedSchema, result.Errors())...)
		}
	}

	return failedReasons
}

func parseValidateError(schema string, resultErrors []gojsonschema.ResultError) []string {
	failedReasons := make([]string, 0)
	for _, desc := range resultErrors {
		// Output string: "demo-v1.(root): url is required"
		failedReasons = append(failedReasons, schema+"."+desc.String())
	}
	return failedReasons
}

func getJSONStr(source string) (string, error) {
	jsonByte, err := httputil.GetByte(source)
	if err != nil {
		return "", err
	}
	buffer := bytes.Buffer{}
	err = json.Compact(&buffer, jsonByte)
	if err != nil {
		return "", err
	}
	return buffer.String(), nil
}

func sendNodeValidationFailedEvent(node *node.Node, failedReasons []string) {
	event.NewNodeValidationFailedPublisher(nats.Client()).Publish(event.NodeValidationFailedData{
		ProfileUrl:    node.ProfileUrl,
		FailedReasons: failedReasons,
		Version:       node.Version,
	})
}
