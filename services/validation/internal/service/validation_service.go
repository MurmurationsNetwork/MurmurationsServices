package service

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/xeipuuv/gojsonschema"

	"github.com/MurmurationsNetwork/MurmurationsServices/common/backoff"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/cryptoutil"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/dateutil"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/event"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/httputil"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/jsonapi"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/nats"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/validatenode"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/validation/config"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/validation/internal/entity"
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

	// Validate against the default schema.
	titles, details, sources, errorStatus := validatenode.ValidateAgainstSchemas(
		config.Conf.Library.InternalURL,
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
		logger.Info(
			"Could not read linked_schemas from profile_url: " + node.ProfileURL,
		)
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

	// Validate against schemes specify inside the profile data.
	titles, details, sources, errorStatus = validatenode.ValidateAgainstSchemas(
		config.Conf.Library.InternalURL,
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
		logger.Info(
			"Could not get JSON string from profile_url: " + node.ProfileURL,
		)
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

func (svc *validationService) sendNodeValidationFailedEvent(
	node *entity.Node,
	FailureReasons *[]jsonapi.Error,
) {
	event.NewNodeValidationFailedPublisher(nats.Client.Client()).
		Publish(event.NodeValidationFailedData{
			ProfileURL:     node.ProfileURL,
			FailureReasons: FailureReasons,
			Version:        node.Version,
		})
}
