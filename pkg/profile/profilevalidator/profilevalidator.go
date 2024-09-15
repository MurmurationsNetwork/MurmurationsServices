package profilevalidator

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

// ProfileValidator is the main struct of this package, which performs
// the validation against the profile JSON.
type ProfileValidator struct {
	// The loader that fetches the profile data.
	ProfileLoader ProfileLoader
	// The Profile JSON data.
	JSON map[string]interface{}
	// The schemas to be used for validation.
	Schemas []string
	// The loader that fetches the schemas.
	SchemaLoader Loader
	// Enable/disable custom validation.
	CustomValidation bool
}

// Validate validates the data against the schemas and returns the validation result.
func (v *ProfileValidator) Validate() *ValidationResult {
	finalResult := NewValidationResult()

	for _, schemaStr := range v.Schemas {
		schema, err := v.SchemaLoader.Load(schemaStr)
		if err != nil {
			finalResult.AppendError(
				"Error loading schema",
				fmt.Sprintf("Error loading schema (%s): %v", schemaStr, err),
				[]string{"pointer", "/linked_schemas"},
				http.StatusNotFound,
			)
			continue
		}

		validationResult, err := schema.Validate(v.ProfileLoader.Load())
		if err != nil {
			finalResult.AppendError(
				"Cannot Validate Document",
				fmt.Sprintf(
					"Error when trying to validate document: %s",
					err.Error(),
				),
				nil,
				http.StatusBadRequest,
			)
			continue
		}

		if !validationResult.Valid() {
			failedTitles, failedDetails, failedSources := parseValidateError(
				schemaStr,
				validationResult.Errors(),
			)
			statusCodes := make([]int, len(failedTitles))
			for i := range statusCodes {
				// Populate it with the same status code for each error.
				statusCodes[i] = http.StatusBadRequest
			}
			finalResult.AppendErrors(
				failedTitles,
				failedDetails,
				failedSources,
				statusCodes,
			)
		}
	}

	if !finalResult.Valid {
		return finalResult
	} else if v.CustomValidation {
		return finalResult.Merge(v.CustomValidate())
	}
	return finalResult
}

// CustomValidate performs custom validation on the JSON data.
func (v *ProfileValidator) CustomValidate() *ValidationResult {
	finalResult := NewValidationResult()

	validators := map[string]CustomValidator{
		"geolocation":      &GeolocationValidator{},
		"name":             &StringValidator{MaxLength: 200, Path: "name"},
		"locality":         &StringValidator{MaxLength: 100, Path: "locality"},
		"region":           &StringValidator{MaxLength: 100, Path: "region"},
		"country_name":     &StringValidator{MaxLength: 100, Path: "country_name"},
		"country_iso_3166": &StringValidator{MaxLength: 2, Path: "country_iso_3166"},
		"primary_url":      &StringValidator{MaxLength: 2000, Path: "primary_url"},
		"tags":             &TagsValidator{},
	}

	for field, validator := range validators {
		if value, exists := v.JSON[field]; exists {
			result := validator.Validate(value)
			if !result.Valid {
				finalResult.Merge(result)
			}
		}
	}

	return finalResult
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
