package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/MurmurationsNetwork/MurmurationsServices/common/backoff"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/cryptoutil"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/dateutil"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/event"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/httputil"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/nats"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/validation/config"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/validation/internal/entity"
	"github.com/xeipuuv/gojsonschema"
)

type ValidationService interface {
	ValidateNode(node *entity.Node)
}

type validationService struct {
}

func NewValidationService() ValidationService {
	return &validationService{}
}

func (svc *validationService) ValidateNode(node *entity.Node) {
	data, err := svc.readFromProfileURL(node.ProfileURL)
	if err != nil {
		logger.Info("Could not read from profile_url: " + node.ProfileURL)
		svc.sendNodeValidationFailedEvent(node, []string{"Could not read from profile_url: " + node.ProfileURL})
		return
	}

	// Validate against the default schema.
	failureReasons := svc.validateAgainstSchemas([]string{"default-v1.0.0"}, node.ProfileURL)
	if len(failureReasons) != 0 {
		logger.Info("Failed to validate against schemas: " + strings.Join(failureReasons, " "))
		svc.sendNodeValidationFailedEvent(node, failureReasons)
		return
	}

	linkedSchemas, ok := getLinkedSchemas(data)
	if !ok {
		logger.Info("Could not read linked_schemas from profile_url: " + node.ProfileURL)
		svc.sendNodeValidationFailedEvent(node, []string{"Could not read linked_schemas from profile_url: " + node.ProfileURL})
		return
	}

	// Validate against schemes specify inside the profile data.
	failureReasons = svc.validateAgainstSchemas(linkedSchemas, node.ProfileURL)
	if len(failureReasons) != 0 {
		logger.Info("Failed to validate against schemas: " + strings.Join(failureReasons, " "))
		svc.sendNodeValidationFailedEvent(node, failureReasons)
		return
	}

	jsonStr, err := getJSONStr(node.ProfileURL)
	if err != nil {
		logger.Info("Could not get JSON string from profile_url: " + node.ProfileURL)
		svc.sendNodeValidationFailedEvent(node, []string{"Could not get JSON string from profile_url: " + node.ProfileURL})
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

func (svc *validationService) readFromProfileURL(profileURL string) (interface{}, error) {
	document := gojsonschema.NewReferenceLoader(profileURL)

	var data interface{}

	// TODO: Need a feature toggle mechanism.
	if true {
		var err error
		data, err = document.LoadJSON()
		if err != nil {
			return nil, err
		}
	} else {
		operation := func() error {
			var err error
			data, err = document.LoadJSON()
			if err != nil {
				return err
			}
			return nil
		}
		err := backoff.NewBackoff(operation, "Could not read from profile_url: "+profileURL)
		if err != nil {
			return nil, err
		}
	}

	return data, nil
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

func (svc *validationService) validateAgainstSchemas(linkedSchemas []string, profileURL string) []string {
	FailureReasons := []string{}

	for _, linkedSchema := range linkedSchemas {
		schemaURL := getSchemaURL(linkedSchema)

		schema, err := gojsonschema.NewSchema(gojsonschema.NewReferenceLoader(schemaURL))
		if err != nil {
			FailureReasons = append(FailureReasons, fmt.Sprintf("Error when trying to read from schema %s: %s", schemaURL, err.Error()))
			continue
		}

		result, err := schema.Validate(gojsonschema.NewReferenceLoader(profileURL))
		if err != nil {
			FailureReasons = append(FailureReasons, "Error when trying to validate document: ", err.Error())
			continue
		}

		if !result.Valid() {
			FailureReasons = append(FailureReasons, svc.parseValidateError(linkedSchema, result.Errors())...)
		}
	}

	return FailureReasons
}

func (svc *validationService) sendNodeValidationFailedEvent(node *entity.Node, FailureReasons []string) {
	event.NewNodeValidationFailedPublisher(nats.Client.Client()).Publish(event.NodeValidationFailedData{
		ProfileURL:     node.ProfileURL,
		FailureReasons: FailureReasons,
		Version:        node.Version,
	})
}

func (svc *validationService) parseValidateError(schema string, resultErrors []gojsonschema.ResultError) []string {
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

func getSchemaURL(linkedSchema string) string {
	return config.Conf.Library.URL + "/schemas/" + linkedSchema + ".json"
}
