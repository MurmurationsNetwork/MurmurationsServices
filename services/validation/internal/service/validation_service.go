package service

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/xeipuuv/gojsonschema"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/cryptoutil"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/dateutil"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/event"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/httputil"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/jsonapi"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/nats"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/retry"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/validatenode"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/validation/config"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/validation/internal/model"
)

type ValidationService interface {
	ValidateNode(node *model.Node)
}

type validationService struct {
}

func NewValidationService() ValidationService {
	return &validationService{}
}

func (svc *validationService) ValidateNode(node *model.Node) {
	data, err := svc.readFromProfileURL(node.ProfileURL)
	if err != nil {
		svc.handleNodeValidationError(node, "Could not read from profile_url")
		return
	}

	// Validate against the default schema. The default schema ensures there is
	// at least one schema defined for validating the node profile.
	titles, details, sources, errorStatus := validatenode.ValidateAgainstSchemas(
		config.Values.Library.InternalURL,
		[]string{"default-v2.0.0"},
		node.ProfileURL,
		"reference",
	)
	if len(titles) != 0 {
		errors := jsonapi.NewError(titles, details, sources, errorStatus)
		logger.Info(
			"Failed to validate against schemas: " + fmt.Sprintf("%v", errors),
		)
		svc.sendNodeValidationFailedEvent(node, &errors)
		return
	}

	linkedSchemas, ok := getLinkedSchemas(data)
	if !ok {
		svc.handleNodeValidationError(
			node,
			"Could not read linked_schemas from profile_url",
		)
		return
	}

	// Validate against the schemas specified in the profile data.
	titles, details, sources, errorStatus = validatenode.ValidateAgainstSchemas(
		config.Values.Library.InternalURL,
		linkedSchemas,
		node.ProfileURL,
		"reference",
	)
	if len(titles) != 0 {
		message := "Failed to validate against schemas: " + strings.Join(
			titles,
			" ",
		)
		logger.Info(message)
		errors := jsonapi.NewError(titles, details, sources, errorStatus)
		svc.sendNodeValidationFailedEvent(node, &errors)
		return
	}

	jsonStr, err := httputil.GetJSONStr(node.ProfileURL)
	if err != nil {
		svc.handleNodeValidationError(
			node,
			"Could not get JSON string from profile_url",
		)
		return
	}

	event.NewNodeValidatedPublisher(nats.Client.Client()).
		Publish(event.NodeValidatedData{
			ProfileURL:  node.ProfileURL,
			ProfileHash: cryptoutil.GetSHA256(jsonStr),
			ProfileStr:  jsonStr,
			LastUpdated: dateutil.GetNowUnix(),
			Version:     node.Version,
		})
}

func (svc *validationService) readFromProfileURL(
	profileURL string,
) (interface{}, error) {
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
		err := retry.Do(operation, "Could not read from profile_url: "+profileURL)
		if err != nil {
			return nil, err
		}
	}

	return data, nil
}

func (svc *validationService) handleNodeValidationError(
	node *model.Node,
	errMsg string,
) {
	logger.Info(fmt.Sprintf("%s: %s", errMsg, node.ProfileURL))
	errors := jsonapi.NewError(
		[]string{"Profile Not Found"},
		[]string{
			fmt.Sprintf(
				"Could not find or read from (invalid JSON) the profile_url: %s",
				node.ProfileURL,
			),
		},
		nil,
		[]int{http.StatusNotFound},
	)
	svc.sendNodeValidationFailedEvent(node, &errors)
}

func (svc *validationService) sendNodeValidationFailedEvent(
	node *model.Node,
	FailureReasons *[]jsonapi.Error,
) {
	event.NewNodeValidationFailedPublisher(nats.Client.Client()).
		Publish(event.NodeValidationFailedData{
			ProfileURL:     node.ProfileURL,
			FailureReasons: FailureReasons,
			Version:        node.Version,
		})
}

func getLinkedSchemas(data interface{}) ([]string, bool) {
	jsonData, ok := data.(map[string]interface{})
	if !ok {
		return nil, false
	}

	linkedSchemasInterface, ok := jsonData["linked_schemas"]
	if !ok {
		return nil, false
	}

	arrInterface, ok := linkedSchemasInterface.([]interface{})
	if !ok {
		return nil, false
	}

	linkedSchemas := make([]string, len(arrInterface))
	for i, data := range arrInterface {
		linkedSchema, ok := data.(string)
		if !ok {
			return nil, false
		}
		linkedSchemas[i] = linkedSchema
	}

	return linkedSchemas, true
}
