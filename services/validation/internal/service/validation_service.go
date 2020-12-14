package service

import (
	"bytes"
	"encoding/json"

	"github.com/MurmurationsNetwork/MurmurationsServices/common/cryptoutil"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/dateutil"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/event"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/httputil"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/nats"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/validation/config"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/validation/internal/domain/node"
	"github.com/xeipuuv/gojsonschema"
)

type ValidationService interface {
	ValidateNode(node *node.Node)
}

type validationService struct {
}

func NewValidationService() ValidationService {
	return &validationService{}
}

func (svc *validationService) ValidateNode(node *node.Node) {
	document := gojsonschema.NewReferenceLoader(node.ProfileURL)
	data, err := document.LoadJSON()
	if err != nil {
		sendNodeValidationFailedEvent(node, []string{"Could not read from profile_url: " + node.ProfileURL})
		return
	}

	// Validate against the default schema.
	failureReasons := validateAgainstSchemas([]string{"default-v1"}, document)
	if len(failureReasons) != 0 {
		sendNodeValidationFailedEvent(node, failureReasons)
		return
	}

	linkedSchemas, ok := getLinkedSchemas(data)
	if !ok {
		sendNodeValidationFailedEvent(node, []string{"Could not read linked_schemas from profile_url: " + node.ProfileURL})
		return
	}

	// Validate against schemes specify inside the profile data.
	failureReasons = validateAgainstSchemas(linkedSchemas, document)
	if len(failureReasons) != 0 {
		sendNodeValidationFailedEvent(node, failureReasons)
		return
	}

	jsonStr, err := getJSONStr(node.ProfileURL)
	if err != nil {
		sendNodeValidationFailedEvent(node, []string{"Could not get JSON string from profile_url: " + node.ProfileURL})
		return
	}

	event.NewNodeValidatedPublisher(nats.Client.Client()).Publish(event.NodeValidatedData{
		ProfileURL:    node.ProfileURL,
		ProfileHash:   cryptoutil.GetSHA256(jsonStr),
		ProfileStr:    jsonStr,
		LastValidated: dateutil.GetNowUnix(),
		Version:       node.Version,
	})
}

func getLinkedSchemas(data interface{}) ([]string, bool) {
	json, ok := data.(map[string]interface{})
	if !ok {
		return nil, false
	}
	_, ok = json["linked_schemas"]
	if !ok {
		return nil, false
	}
	arrInterface, ok := json["linked_schemas"].([]interface{})
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
	FailureReasons := []string{}

	for _, linkedSchema := range linkedSchemas {
		schemaURL := getSchemaURL(linkedSchema)

		schema, err := gojsonschema.NewSchema(gojsonschema.NewReferenceLoader(schemaURL))
		if err != nil {
			FailureReasons = append(FailureReasons, "Could not read from schema: "+schemaURL)
			continue
		}

		result, err := schema.Validate(document)
		if err != nil {
			FailureReasons = append(FailureReasons, "Error when trying to validate document: ", err.Error())
			continue
		}

		if !result.Valid() {
			FailureReasons = append(FailureReasons, parseValidateError(linkedSchema, result.Errors())...)
		}
	}

	return FailureReasons
}

func parseValidateError(schema string, resultErrors []gojsonschema.ResultError) []string {
	FailureReasons := make([]string, 0)
	for _, desc := range resultErrors {
		// Output string: "demo-v1.(root): url is required"
		FailureReasons = append(FailureReasons, schema+"."+desc.String())
	}
	return FailureReasons
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

func sendNodeValidationFailedEvent(node *node.Node, FailureReasons []string) {
	event.NewNodeValidationFailedPublisher(nats.Client.Client()).Publish(event.NodeValidationFailedData{
		ProfileURL:     node.ProfileURL,
		FailureReasons: FailureReasons,
		Version:        node.Version,
	})
}

func getSchemaURL(linkedSchema string) string {
	return config.Conf.Library.URL + "/schemas/" + linkedSchema + ".json"
}
