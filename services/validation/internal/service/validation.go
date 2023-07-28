package service

import (
	"fmt"
	"net/http"

	"github.com/xeipuuv/gojsonschema"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/cryptoutil"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/dateutil"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/event"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/httputil"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/jsonapi"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/jsonutil"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/nats"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/retry"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/schemavalidator"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/validation/config"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/validation/internal/model"
)

const DefaultSchema = "default-v2.0.0"

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
		logger.Info(
			"Failed to read from profile URL: " + fmt.Sprintf("%v", errors),
		)
		svc.sendNodeValidationFailedEvent(node, &errors)
		return
	}

	// Validate against the default schema. The default schema ensures there is
	// at least one schema defined for validating the node profile.
	validator, err := schemavalidator.NewBuilder().
		WithURLSchemas(config.Values.Library.InternalURL, []string{DefaultSchema}).
		WithURLProfile(node.ProfileURL).
		Build()
	if err != nil {
		// Log the error for internal debugging and auditing.
		logger.Error("Failed to build schema validator", err)

		errors := jsonapi.NewError(
			[]string{"Internal Server Error"},
			[]string{
				"An error occurred while validating the profile data. Please try again later.",
			},
			nil,
			[]int{http.StatusInternalServerError},
		)
		svc.sendNodeValidationFailedEvent(node, &errors)
		return
	}

	result := validator.Validate()
	if !result.Valid {
		errors := jsonapi.NewError(
			result.ErrorMessages,
			result.Details,
			result.Sources,
			result.ErrorStatus,
		)
		svc.sendNodeValidationFailedEvent(node, &errors)
		return
	}

	linkedSchemas, ok := getLinkedSchemas(data)
	if !ok {
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
		return
	}

	// Validate against the schemas specified in the profile data.
	validator, err = schemavalidator.NewBuilder().
		WithURLSchemas(config.Values.Library.InternalURL, linkedSchemas).
		WithURLProfile(node.ProfileURL).
		Build()
	if err != nil {
		// Log the error for internal debugging and auditing.
		logger.Error("Failed to build schema validator", err)

		errors := jsonapi.NewError(
			[]string{"Internal Server Error"},
			[]string{
				"An error occurred while validating the profile data. Please try again later.",
			},
			nil,
			[]int{http.StatusInternalServerError},
		)
		svc.sendNodeValidationFailedEvent(node, &errors)
		return
	}

	result = validator.Validate()
	if !result.Valid {
		errors := jsonapi.NewError(
			result.ErrorMessages,
			result.Details,
			result.Sources,
			result.ErrorStatus,
		)
		svc.sendNodeValidationFailedEvent(node, &errors)
		return
	}

	jsonStr, err := httputil.GetJSONStr(node.ProfileURL)
	if err != nil {
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
		return
	}

	// Normalize the primary URL.
	profileJSON := jsonutil.ToJSON(jsonStr)
	if profileJSON["primary_url"] != nil {
		normalizedURL, err := NormalizeURL(profileJSON["primary_url"].(string))
		if err != nil {
			errors := jsonapi.NewError(
				[]string{"Primary URL Validation Failed"},
				[]string{
					fmt.Sprintf(
						"The primary URL is invalid: %s.",
						profileJSON["primary_url"].(string),
					),
				},
				nil,
				[]int{http.StatusBadRequest},
			)
			svc.sendNodeValidationFailedEvent(node, &errors)
			return
		}
		profileJSON["primary_url"] = normalizedURL
	}

	event.NewNodeValidatedPublisher(nats.Client.Client()).
		Publish(event.NodeValidatedData{
			ProfileURL:  node.ProfileURL,
			ProfileHash: cryptoutil.GetSHA256(jsonStr),
			// Provides the updated version of the profile for later use.
			ProfileStr:  jsonutil.ToString(profileJSON),
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
		err := retry.Do(operation)
		if err != nil {
			return nil, err
		}
	}

	return data, nil
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
