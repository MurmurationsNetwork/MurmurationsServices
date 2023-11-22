package service

import (
	"fmt"
	"net/http"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/dateutil"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/event"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/httputil"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/jsonapi"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/jsonutil"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/nats"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/profile/profilehasher"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/profile/profilevalidator"
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
	profileStr, err := httputil.GetJSONStr(node.ProfileURL)
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

	if err := svc.validateAgainstDefaultSchema(profileStr, node); err != nil {
		return
	}
	if err := svc.validateAgainstLinkedSchemas(profileStr, node); err != nil {
		return
	}

	profileHash, err := profilehasher.NewFromString(profileStr, config.Values.Library.InternalURL).
		Hash()
	if err != nil {
		logger.Error("Failed to generate a hash for the profile_url: ", err)
		errors := jsonapi.NewError(
			[]string{"Profile Hashing Failed"},
			[]string{
				fmt.Sprintf(
					"Failed to generate a hash for the profile_url: %s. Please try again later.",
					node.ProfileURL,
				),
			},
			nil,
			[]int{http.StatusInternalServerError},
		)
		svc.sendNodeValidationFailedEvent(node, &errors)
		return
	}

	updatedProfileJSON := jsonutil.ToJSON(profileStr)
	if updatedProfileJSON["primary_url"] != nil {
		normalizedURL, err := NormalizeURL(
			updatedProfileJSON["primary_url"].(string),
		)
		if err != nil {
			errors := jsonapi.NewError(
				[]string{"Primary URL Validation Failed"},
				[]string{
					fmt.Sprintf(
						"The primary URL is invalid: %s.",
						updatedProfileJSON["primary_url"].(string),
					),
				},
				nil,
				[]int{http.StatusBadRequest},
			)
			svc.sendNodeValidationFailedEvent(node, &errors)
			return
		}
		updatedProfileJSON["primary_url"] = normalizedURL
	}

	event.NewNodeValidatedPublisher(nats.Client.Client()).
		Publish(event.NodeValidatedData{
			ProfileURL:  node.ProfileURL,
			ProfileHash: profileHash,
			// Provides the updated version of the profile for later use.
			ProfileStr:  jsonutil.ToString(updatedProfileJSON),
			LastUpdated: dateutil.GetNowUnix(),
			Version:     node.Version,
		})
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

// validateAgainstDefaultSchema handles the validation of the node's profile against the default schema.
func (svc *validationService) validateAgainstDefaultSchema(
	profileStr string,
	node *model.Node,
) error {
	validator, err := profilevalidator.NewBuilder().
		WithStrProfile(profileStr).
		WithURLSchemas(config.Values.Library.InternalURL, []string{DefaultSchema}).
		Build()
	if err != nil {
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
		return err
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
		return fmt.Errorf("validation failed")
	}

	return nil
}

// validateAgainstLinkedSchemas handles the extraction and validation of the node's profile against linked schemas.
func (svc *validationService) validateAgainstLinkedSchemas(
	profileStr string,
	node *model.Node,
) error {
	linkedSchemas, err := getLinkedSchemas(profileStr)
	if err != nil {
		errors := jsonapi.NewError(
			[]string{"Profile Validation Error"},
			[]string{err.Error()},
			nil,
			[]int{http.StatusBadRequest},
		)
		svc.sendNodeValidationFailedEvent(node, &errors)
		return err
	}

	validator, err := profilevalidator.NewBuilder().
		WithStrProfile(profileStr).
		WithURLSchemas(config.Values.Library.InternalURL, linkedSchemas).
		WithCustomValidation().
		Build()
	if err != nil {
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
		return err
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
		return fmt.Errorf("validation against linked schemas failed")
	}

	return nil
}

func getLinkedSchemas(profileStr string) ([]string, error) {
	jsonData := jsonutil.ToJSON(profileStr)

	linkedSchemasInterface, ok := jsonData["linked_schemas"]
	if !ok {
		return nil, fmt.Errorf("linked schemas not found in profile")
	}

	arrInterface, ok := linkedSchemasInterface.([]interface{})
	if !ok {
		return nil, fmt.Errorf("linked schemas is not an array")
	}

	if len(arrInterface) == 0 {
		return nil, fmt.Errorf("empty linked schemas array")
	}

	linkedSchemas := make([]string, len(arrInterface))
	for i, data := range arrInterface {
		linkedSchema, ok := data.(string)
		if !ok {
			return nil, fmt.Errorf("invalid schema type in linked schemas")
		}
		linkedSchemas[i] = linkedSchema
	}

	return linkedSchemas, nil
}
