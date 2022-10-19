package validatenode

import (
	"fmt"
	"github.com/xeipuuv/gojsonschema"
	"net/http"
	"strings"
)

func parseValidateError(schema string, resultErrors []gojsonschema.ResultError) ([]string, []string, [][]string) {
	var (
		failedTitles, failedDetails []string
		failedSources               [][]string
	)
	for _, desc := range resultErrors {
		// title
		failedType := desc.Type()

		// details
		var expected, given, min, max, property, failedDetail, failedField string
		for index, value := range desc.Details() {
			if index == "expected" {
				expected = value.(string)
			} else if index == "given" {
				given = value.(string)
			} else if index == "min" {
				min = fmt.Sprint(value)
			} else if index == "max" {
				max = fmt.Sprint(value)
			} else if index == "property" {
				property = value.(string)
			}
		}

		if failedType == "invalid_type" {
			failedType = "Invalid Type"
			failedDetail = "Expected: " + expected + " - Given: " + given + " - Schema: " + schema
		} else if failedType == "number_gte" {
			failedType = "Invalid Amount"
			failedDetail = "Amount must be greater than or equal to " + min + " - Schema: " + schema
		} else if failedType == "number_lte" {
			failedType = "Invalid Amount"
			failedDetail = "Amount must be less than or equal to " + max + " - Schema: " + schema
		} else if failedType == "required" {
			failedType = "Missing Required Property"
			if desc.Field() == "(root)" {
				failedDetail = "The `/" + property + "` property is required."
			} else {
				failedDetail = "The `/" + desc.Field() + "/" + property + "` property is required."
			}
		}

		// append title and detail
		failedTitles = append(failedTitles, failedType)
		failedDetails = append(failedDetails, failedDetail)

		// sources
		if desc.Field() == "(root)" && property != "" {
			failedField = "/" + property
		} else if property != "" {
			failedField = "/" + strings.Replace(desc.Field(), ".", "/", -1) + "/" + property
		} else {
			failedField = "/" + strings.Replace(desc.Field(), ".", "/", -1)
		}
		failedSources = append(failedSources, []string{"pointer", failedField})
	}

	return failedTitles, failedDetails, failedSources
}

func getSchemaURL(schemaUrl string, linkedSchema string) string {
	return schemaUrl + "/v1/schema/" + linkedSchema
}

func ValidateAgainstSchemas(schemaUrl string, linkedSchemas []string, validateData string, schemaLoader string) ([]string, []string, [][]string, []int) {
	var (
		titles, details []string
		sources         [][]string
		errorStatus     []int
	)

	for _, linkedSchema := range linkedSchemas {
		schemaURL := getSchemaURL(schemaUrl, linkedSchema)

		schema, err := gojsonschema.NewSchema(gojsonschema.NewReferenceLoader(schemaURL))
		if err != nil {
			titles = append(titles, []string{"Schema Not Found"}...)
			details = append(details, []string{"Could not locate the following schema in the library: " + linkedSchema}...)
			sources = append(sources, [][]string{{"pointer", "/linked_schemas"}}...)
			errorStatus = append(errorStatus, http.StatusNotFound)
			continue
		}

		var result *gojsonschema.Result
		if schemaLoader == "reference" {
			result, err = schema.Validate(gojsonschema.NewReferenceLoader(validateData))
		} else {
			result, err = schema.Validate(gojsonschema.NewStringLoader(validateData))
		}
		if err != nil {
			titles = append(titles, "Cannot Validate Document")
			details = append(details, []string{"Error when trying to validate document: ", err.Error()}...)
			errorStatus = append(errorStatus, http.StatusBadRequest)
			continue
		}

		if !result.Valid() {
			failedTitles, failedDetails, failedSources := parseValidateError(linkedSchema, result.Errors())
			titles = append(titles, failedTitles...)
			details = append(details, failedDetails...)
			sources = append(sources, failedSources...)
			for i := 0; i < len(titles); i++ {
				errorStatus = append(errorStatus, http.StatusBadRequest)
			}
		}
	}

	return titles, details, sources, errorStatus
}
