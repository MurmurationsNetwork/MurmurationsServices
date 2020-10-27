package service

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/common/dateutil"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/event"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/validation/internal/datasource/nats"
	"github.com/xeipuuv/gojsonschema"
)

var (
	ValidationService validationServiceInterface = &validationService{}
)

type validationServiceInterface interface {
	ValidateNode(profileUrl string, linkedSchemas []string)
}

type validationService struct{}

func (v *validationService) ValidateNode(profileUrl string, linkedSchemas []string) {
	// DISCUSS (2020/10/27): When will we have multiple schemas against a single profile url?
	// 						 They only provide one schema validate multiple documents.
	document := gojsonschema.NewReferenceLoader(profileUrl)

	for _, schemaURL := range linkedSchemas {
		schema, err := gojsonschema.NewSchema(gojsonschema.NewReferenceLoader(schemaURL))
		if err != nil {
			sendNodeValidationFailedEvent(profileUrl, []string{"Could not read from schema: " + schemaURL})
			return
		}

		result, err := schema.Validate(document)
		if err != nil {
			sendNodeValidationFailedEvent(profileUrl, []string{"Could not read from profile url: " + profileUrl})
			return
		}

		if !result.Valid() {
			failedReasons := parseResultError(result.Errors())
			sendNodeValidationFailedEvent(profileUrl, failedReasons)
			return
		}
	}

	event.NewNodeValidatedPublisher(nats.Client()).Publish(event.NodeValidatedData{
		ProfileUrl:    profileUrl,
		LastValidated: dateutil.GetNowUnix(),
	})
}

func parseResultError(resultErrors []gojsonschema.ResultError) []string {
	failedReasons := make([]string, 0)
	for _, desc := range resultErrors {
		failedReasons = append(failedReasons, desc.String())
	}
	return failedReasons
}

func sendNodeValidationFailedEvent(profileUrl string, failedReasons []string) {
	event.NewNodeValidationFailedPublisher(nats.Client()).Publish(event.NodeValidationFailedData{
		ProfileUrl:    profileUrl,
		FailedReasons: failedReasons,
	})
}
