package services

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/common/date_utils"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/events"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/validation/internal/datasources/nats"
	"github.com/xeipuuv/gojsonschema"
)

var (
	ValidationService validationServiceInterface = &validationService{}
)

type validationServiceInterface interface {
	ValidateNode(profileUrl string, linkedSchemas []string) error
}

type validationService struct{}

func (v *validationService) ValidateNode(profileUrl string, linkedSchemas []string) error {
	// DISCUSS (2020/10/27): When will we have multiple schemas against a single profile url?
	// 						 They only provide one schema validate multiple documents.
	document := gojsonschema.NewReferenceLoader(profileUrl)

	for _, schemaURL := range linkedSchemas {
		schema, err := gojsonschema.NewSchema(gojsonschema.NewReferenceLoader(schemaURL))
		if err != nil {
			sendNodeValidationFailedEvent(profileUrl, []string{"Could not read from schema: " + schemaURL})
			return err
		}

		result, err := schema.Validate(document)
		if err != nil {
			sendNodeValidationFailedEvent(profileUrl, []string{"Could not read from profile url: " + profileUrl})
			return err
		}

		if !result.Valid() {
			failedReasons := parseResultError(result.Errors())
			sendNodeValidationFailedEvent(profileUrl, failedReasons)
			return err
		}
	}

	events.NewNodeValidatedPublisher(nats.Client()).Publish(events.NodeValidatedData{
		ProfileUrl:    profileUrl,
		LastValidated: date_utils.GetNowUnix(),
	})

	return nil
}

func parseResultError(resultErrors []gojsonschema.ResultError) []string {
	failedReasons := make([]string, 0)
	for _, desc := range resultErrors {
		failedReasons = append(failedReasons, desc.String())
	}
	return failedReasons
}

func sendNodeValidationFailedEvent(profileUrl string, failedReasons []string) {
	events.NewNodeValidationFailedPublisher(nats.Client()).Publish(events.NodeValidationFailedData{
		ProfileUrl:    profileUrl,
		FailedReasons: failedReasons,
	})
}
