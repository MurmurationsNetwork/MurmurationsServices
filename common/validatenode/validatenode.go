package validatenode

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/xeipuuv/gojsonschema"
)

func parseValidateError(
	schema string,
	resultErrors []gojsonschema.ResultError,
) ([]string, []string, [][]string) {
	var (
		failedTitles, failedDetails []string
		failedSources               [][]string
	)
	for _, desc := range resultErrors {
		// title
		failedType := desc.Type()

		// details
		var expected, given, min, max, property, pattern, failedDetail, failedField string
		for index, value := range desc.Details() {
			switch index {
			case "expected":
				expected = value.(string)
			case "given":
				given = value.(string)
			case "min":
				min = fmt.Sprint(value)
			case "max":
				max = fmt.Sprint(value)
			case "property":
				property = value.(string)
			case "pattern":
				pattern = fmt.Sprint(value)
			}
		}

		switch failedType {
		case "invalid_type":
			failedType = "Invalid Type"
			failedDetail = "Expected: " + expected + " - Given: " + given + " - Schema: " + schema
		case "number_gte":
			failedType = "Invalid Amount"
			failedDetail = "Amount must be greater than or equal to " + min + " - Schema: " + schema
		case "number_lte":
			failedType = "Invalid Amount"
			failedDetail = "Amount must be less than or equal to " + max + " - Schema: " + schema
		case "required":
			failedType = "Missing Required Property"
			if desc.Field() == "(root)" {
				failedDetail = "The `" + property + "` property is required - Schema: " + schema
			} else {
				failedDetail = "The `" + desc.Field() + "/" + property + "` property is required - Schema: " + schema
			}
		case "array_min_items":
			failedType = "Not Enough Items"
			failedDetail = "There are not enough items in the array - Schema: " + schema
		case "array_max_items":
			failedType = "Too Many Items"
			failedDetail = "There are too many items in the array - Schema: " + schema
		case "pattern":
			failedType = "Pattern Mismatch"
			failedDetail = "The submitted data does not match the required pattern: '" + pattern + "' - Schema: " + schema
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
	return schemaUrl + "/v2/schemas/" + linkedSchema
}

func ValidateAgainstSchemas(
	schemaUrl string,
	linkedSchemas []string,
	validateData string,
	schemaLoader string,
) ([]string, []string, [][]string, []int) {
	var (
		titles, details []string
		sources         [][]string
		errorStatus     []int
	)

	for _, linkedSchema := range linkedSchemas {
		schemaURL := getSchemaURL(schemaUrl, linkedSchema)

		schema, err := gojsonschema.NewSchema(
			gojsonschema.NewReferenceLoader(schemaURL),
		)
		if err != nil {
			titles = append(titles, []string{"Schema Not Found"}...)
			details = append(
				details,
				[]string{
					"Could not locate the following schema in the Library: " + linkedSchema,
				}...)
			sources = append(
				sources,
				[][]string{{"pointer", "/linked_schemas"}}...)
			errorStatus = append(errorStatus, http.StatusNotFound)
			continue
		}

		var result *gojsonschema.Result
		if schemaLoader == "reference" {
			result, err = schema.Validate(
				gojsonschema.NewReferenceLoader(validateData),
			)
		} else {
			result, err = schema.Validate(gojsonschema.NewStringLoader(validateData))
		}
		if err != nil {
			titles = append(titles, "Cannot Validate Document")
			details = append(
				details,
				[]string{
					"Error when trying to validate document: ",
					err.Error(),
				}...)
			errorStatus = append(errorStatus, http.StatusBadRequest)
			continue
		}

		if !result.Valid() {
			failedTitles, failedDetails, failedSources := parseValidateError(
				linkedSchema,
				result.Errors(),
			)
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

func ValidateAgainstSchemasWithoutURL(
	linkedSchemas []string,
	validateSchemas []string,
	validateData map[string]interface{},
) ([]string, []string, [][]string, []int) {
	var (
		titles, details []string
		sources         [][]string
		errorStatus     []int
	)

	for i, linkedSchema := range linkedSchemas {
		schema, err := gojsonschema.NewSchema(
			gojsonschema.NewStringLoader(linkedSchema),
		)
		if err != nil {
			titles = append(titles, []string{"Schema Not Found"}...)
			details = append(
				details,
				[]string{
					"Could not parse the schema: " + validateSchemas[i],
				}...)
			sources = append(
				sources,
				[][]string{{"pointer", "/linked_schemas"}}...)
			errorStatus = append(errorStatus, http.StatusNotFound)
			continue
		}

		var result *gojsonschema.Result
		result, err = schema.Validate(gojsonschema.NewGoLoader(validateData))
		if err != nil {
			titles = append(titles, "Cannot Validate Document")
			details = append(
				details,
				[]string{
					"Error when trying to validate document: ",
					err.Error(),
				}...)
			errorStatus = append(errorStatus, http.StatusBadRequest)
			continue
		}

		if !result.Valid() {
			failedTitles, failedDetails, failedSources := parseValidateError(
				validateSchemas[i],
				result.Errors(),
			)
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
