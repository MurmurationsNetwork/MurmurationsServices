package schemavalidator

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
