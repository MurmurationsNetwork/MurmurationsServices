package schemavalidator

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/xeipuuv/gojsonschema"
)

// Loader is the interface that wraps the Load method.
type Loader interface {
	// Load fetches the JSON schema from a source and returns it.
	Load(string) (*gojsonschema.Schema, error)
}

// URLSchemaLoader is a schema loader that loads schema from a URL.
type URLSchemaLoader struct {
	// The base URL for the schemas.
	BaseURL string
}

// Load implements the Loader interface.
func (ul *URLSchemaLoader) Load(
	linkedSchema string,
) (*gojsonschema.Schema, error) {
	schemaURL := getSchemaURL(ul.BaseURL, linkedSchema)
	return gojsonschema.NewSchema(gojsonschema.NewReferenceLoader(schemaURL))
}

// StrSchemaLoader is a schema loader that loads schema from a string.
type StrSchemaLoader struct{}

// Load implements the Loader interface.
func (sl *StrSchemaLoader) Load(
	source string,
) (*gojsonschema.Schema, error) {
	return gojsonschema.NewSchema(gojsonschema.NewStringLoader(source))
}

// ProfileLoader is the interface that wraps the Load method.
type ProfileLoader interface {
	// Load fetches the data to be validated and returns it.
	Load() gojsonschema.JSONLoader
}

// URLProfileLoader is a profile loader that loads data from a URL.
type URLProfileLoader struct {
	// The URL of the data.
	dataURL string
}

// Load implements the ProfileLoader interface.
func (rv *URLProfileLoader) Load() gojsonschema.JSONLoader {
	return gojsonschema.NewReferenceLoader(rv.dataURL)
}

// StrProfileLoader is a profile loader that loads data from a string.
type StrProfileLoader struct {
	// The string of the data.
	dataString string
}

// Load implements the ProfileLoader interface.
func (sv *StrProfileLoader) Load() gojsonschema.JSONLoader {
	return gojsonschema.NewStringLoader(sv.dataString)
}

// MapProfileLoader is a profile loader that loads data from a map.
type MapProfileLoader struct {
	// The map of the data.
	dataMap map[string]interface{}
}

// Load implements the ProfileLoader interface.
func (mv *MapProfileLoader) Load() gojsonschema.JSONLoader {
	return gojsonschema.NewGoLoader(mv.dataMap)
}

// SchemaValidator is the main struct of this package, which performs
// the validation against the JSON schemas.
type SchemaValidator struct {
	// The schemas to be used for validation.
	Schemas []string
	// The loader that fetches the schemas.
	SchemaLoader Loader
	// The loader that fetches the profile data.
	ProfileLoader ProfileLoader
}

// Validate validates the data against the schemas and returns the validation result.
func (v *SchemaValidator) Validate() *ValidationResult {
	var (
		errorMessages []string
		details       []string
		sources       [][]string
		errorStatus   []int
	)

	for _, schemaStr := range v.Schemas {
		schema, err := v.SchemaLoader.Load(schemaStr)
		if err != nil {
			errorMessages = append(errorMessages, "Error loading schema")
			details = append(
				details,
				fmt.Sprintf(
					"Error loading schema (%s): %v",
					schemaStr,
					err,
				),
			)
			sources = append(sources, []string{"pointer", "/linked_schemas"})
			errorStatus = append(errorStatus, http.StatusNotFound)
			continue
		}

		result, err := schema.Validate(v.ProfileLoader.Load())
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
				schemaStr,
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

// getSchemaURL constructs the full schema URL and returns it.
func getSchemaURL(libraryURL string, linkedSchema string) string {
	return fmt.Sprintf("%s/v2/schemas/%s", libraryURL, linkedSchema)
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
