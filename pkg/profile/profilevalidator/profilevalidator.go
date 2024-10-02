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

// ProfileValidator is responsible for validating profile JSON data against schemas.
type ProfileValidator struct {
	ProfileLoader    ProfileLoader          // Fetches the profile data.
	ProfileJSON      map[string]interface{} // The loaded profile JSON data.
	SchemaNames      []string               // Names of the schemas (for context or errors).
	SchemaReferences []string               // For URL-based schemas (used to fetch schema content).
	LoadedSchemas    []string               // For JSON-based schemas (actual content).
	SchemaLoader     Loader                 // Loader for fetching schema content.
}

// Validate performs validation of the profile JSON against the provided schemas
// and returns the aggregated validation result.
func (v *ProfileValidator) Validate() *ValidationResult {
	finalResult := NewValidationResult()

	// Use either loaded JSON schemas or schema references (both are set during initialization).
	schemasToValidate := v.LoadedSchemas
	if len(schemasToValidate) == 0 {
		schemasToValidate = v.SchemaReferences
	}

	// Iterate over each schema for validation.
	for i, schema := range schemasToValidate {
		// Load the schema using the SchemaLoader.
		loadedSchema, err := v.SchemaLoader.Load(schema)
		if err != nil {
			finalResult.AppendError(
				"Error loading schema",
				fmt.Sprintf("Error loading schema (%s): %v", v.SchemaNames[i], err),
				[]string{"pointer", "/linked_schemas"},
				http.StatusNotFound,
			)
			continue
		}

		// Validate the profile JSON against the loaded schema.
		validationResult, err := loadedSchema.Validate(v.ProfileLoader.Load())
		if err != nil {
			finalResult.AppendError(
				"Cannot Validate Document",
				fmt.Sprintf("Error validating document: %s", err.Error()),
				nil,
				http.StatusBadRequest,
			)
			continue
		}

		// If validation fails, collect and append the errors.
		if !validationResult.Valid() {
			titles, details, sources := parseValidateError(v.SchemaNames[i], validationResult.Errors())

			// Assign the same status code for each validation error.
			statusCodes := make([]int, len(titles))
			for i := range statusCodes {
				statusCodes[i] = http.StatusBadRequest
			}

			// Append all validation errors to the final result.
			finalResult.AppendErrors(titles, details, sources, statusCodes)
		}
	}

	return finalResult
}

// getSchemaURL constructs the full schema URL and returns it.
func getSchemaURL(libraryURL string, linkedSchema string) string {
	return fmt.Sprintf("%s/v2/schemas/%s", libraryURL, linkedSchema)
}

func parseValidateError(
	schemaName string,
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
			failedDetail = "Expected: " + expected + " - Given: " + given + " - Schema: " + schemaName
		case "number_gte":
			failedType = "Invalid Amount"
			failedDetail = "Amount must be greater than or equal to " + min + " - Schema: " + schemaName
		case "number_lte":
			failedType = "Invalid Amount"
			failedDetail = "Amount must be less than or equal to " + max + " - Schema: " + schemaName
		case "required":
			failedType = "Missing Required Property"
			if desc.Field() == "(root)" {
				failedDetail = "The `" + property + "` property is required - Schema: " + schemaName
			} else {
				failedDetail = "The `" + desc.Field() + "/" + property + "` property is required - Schema: " + schemaName
			}
		case "array_min_items":
			failedType = "Not Enough Items"
			failedDetail = "There are not enough items in the array - Minimum is " + min + " - Schema: " + schemaName
		case "array_max_items":
			failedType = "Too Many Items"
			failedDetail = "There are too many items in the array - Maximum is " + max + " - Schema: " + schemaName
		case "pattern":
			failedType = "Pattern Mismatch"
			failedDetail = "The submitted data does not match the required pattern: '" + pattern + "' - Schema: " + schemaName
		case "enum":
			failedType = "Invalid Value"
			failedDetail = "The submitted data is not a valid value from the list of allowed values - Schema: " + schemaName
		case "unique":
			failedType = "Duplicate Value"
			failedDetail = "The submitted data contains a duplicate value - Schema: " + schemaName
		case "string_lte":
			failedType = "Invalid Length"
			failedDetail = "Amount must be less than or equal to " + max + " - Schema: " + schemaName
		case "string_gte":
			failedType = "Invalid Length"
			failedDetail = "Amount must be greater than or equal to " + min + " - Schema: " + schemaName
		// condition_else and condition_then are not errors, they are conditions - no need to report them
		case "condition_else":
			continue
		case "condition_then":
			continue
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
