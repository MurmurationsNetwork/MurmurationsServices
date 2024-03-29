package index

import "fmt"

// DatabaseError struct represents a custom error type for database operations.
type DatabaseError struct {
	// Human readable error message.
	Message string
	// Wrapped error.
	Err error
}

// Error conforms to go conventions.
func (e DatabaseError) Error() string {
	return fmt.Sprintf(
		"Encountered a DB error. %s: %s",
		e.Message,
		e.Err,
	)
}

// Unwrap conforms to go conventions.
func (e DatabaseError) Unwrap() error {
	return e.Err
}

// NotFoundError struct represents a custom error type for not found situations.
type NotFoundError struct {
	// Wrapped error
	Err error
}

// Error function to represent the NotFoundError as a string.
func (e NotFoundError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("Record not found: %v", e.Err)
	}
	return "Record not found"
}

// Unwrap conforms to go conventions.
func (e NotFoundError) Unwrap() error {
	return e.Err
}

// ValidationError struct represents a custom error type for validation failures.
type ValidationError struct {
	// Field that failed validation.
	Field string
	// Reason why validation failed.
	Reason string
}

// Error conforms to go conventions.
func (e ValidationError) Error() string {
	return fmt.Sprintf("Validation failed on field '%s': %s", e.Field, e.Reason)
}

const (
	// HTTP request failure.
	ErrorHTTPRequestFailed = 1
	// Profile URL check failure.
	ErrorProfileURLCheckFail = 2
	// Profile URL still exists.
	ErrorProfileStillExists = 3
)

// DeleteNodeError defines errors during node deletion, including context.
type DeleteNodeError struct {
	Message    string // General error message
	Detail     string // Error details
	NodeID     string // Affected node ID
	ProfileURL string // Associated profile URL
	ErrorCode  int    // Error categorization code
}

// Error conforms to go conventions.
func (e DeleteNodeError) Error() string {
	return fmt.Sprintf(
		"Message: %s, Detail: %s, Profile URL: %s, Node ID: %s",
		e.Message,
		e.Detail,
		e.ProfileURL,
		e.NodeID,
	)
}
