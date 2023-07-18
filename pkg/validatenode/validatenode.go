package validatenode

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/xeipuuv/gojsonschema"
)

// ValidationResult is the results from a schema validation operation.
type ValidationResult struct {
	// Whether the validation passed.
	Valid bool
	// Titles of the validation errors.
	ErrorMessages []string
	// Detailed descriptions of the validation errors.
	Details []string
	// Sources indicates the failing pieces of data.
	Sources [][]string
	// HTTP status codes associated with each error.
	ErrorStatus []int
}

// ValidateAgainstSchemas validates a data string against a set of schemas.
func ValidateAgainstSchemas(
	schemaURL string,
	linkedSchemas []string,
	validateData string,
	schemaLoader string,
) *ValidationResult {
	var (
		errorMessages []string
		details       []string
		sources       [][]string
		errorStatus   []int
	)

	for _, linkedSchema := range linkedSchemas {
		// Construct the full schema URL.
		schemaURL := getSchemaURL(schemaURL, linkedSchema)

		// Load the schema.
		schema, err := gojsonschema.NewSchema(
			gojsonschema.NewReferenceLoader(schemaURL),
		)
		if err != nil {
			errorMessages = append(errorMessages, "Error loading schema")
			details = append(
				details,
				fmt.Sprintf(
					"Error loading schema (%s): %v",
					linkedSchema,
					err,
				),
			)
			sources = append(sources, []string{"pointer", "/linked_schemas"})
			errorStatus = append(errorStatus, http.StatusNotFound)
			continue
		}

		var result *gojsonschema.Result

		// Validate the data against the schema.
		if schemaLoader == "reference" {
			result, err = schema.Validate(
				gojsonschema.NewReferenceLoader(validateData),
			)
		} else {
			result, err = schema.Validate(gojsonschema.NewStringLoader(validateData))
		}

		if err != nil {
			errorMessages = append(errorMessages, "Cannot Validate Document")
			details = append(
				details,
				fmt.Sprintf(
					"Error when trying to validate document: %s",
					err.Error(),
				),
			)
			errorStatus = append(errorStatus, http.StatusBadRequest)
			continue
		}

		if !result.Valid() {
			failedTitles, failedDetails, failedSources := parseValidateError(
				linkedSchema,
				result.Errors(),
			)
			errorMessages = append(errorMessages, failedTitles...)
			details = append(details, failedDetails...)
			sources = append(sources, failedSources...)
			for i := 0; i < len(failedTitles); i++ {
				errorStatus = append(errorStatus, http.StatusBadRequest)
			}
		}
	}

	return &ValidationResult{
		Valid:         len(errorMessages) == 0,
		ErrorMessages: errorMessages,
		Details:       details,
		Sources:       sources,
		ErrorStatus:   errorStatus,
	}
}

func ValidateAgainstSchemasWithoutURL(
	linkedSchemas []string,
	validateSchemas []string,
	validateData map[string]interface{},
) *ValidationResult {
	var (
		errorMessages, details []string
		sources                [][]string
		errorStatus            []int
	)

	for i, linkedSchema := range linkedSchemas {
		schema, err := gojsonschema.NewSchema(
			gojsonschema.NewStringLoader(linkedSchema),
		)
		if err != nil {
			errorMessages = append(
				errorMessages,
				[]string{"Schema Not Found"}...)
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
			errorMessages = append(errorMessages, "Cannot Validate Document")
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
			failederrorMessages, failedDetails, failedSources := parseValidateError(
				validateSchemas[i],
				result.Errors(),
			)
			errorMessages = append(errorMessages, failederrorMessages...)
			details = append(details, failedDetails...)
			sources = append(sources, failedSources...)
			for i := 0; i < len(errorMessages); i++ {
				errorStatus = append(errorStatus, http.StatusBadRequest)
			}
		}
	}

	return &ValidationResult{
		Valid:         len(errorMessages) == 0,
		ErrorMessages: errorMessages,
		Details:       details,
		Sources:       sources,
		ErrorStatus:   errorStatus,
	}
}

func parseValidateError(
	schema string,
	resultErrors []gojsonschema.ResultError,
) ([]string, []string, [][]string) {
	failedTitles := make([]string, 0, len(resultErrors))
	failedDetails := make([]string, 0, len(resultErrors))
	failedSources := make([][]string, 0, len(resultErrors))

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

func getSchemaURL(schemaURL string, linkedSchema string) string {
	return schemaURL + "/v2/schemas/" + linkedSchema
}
