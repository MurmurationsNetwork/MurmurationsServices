package library

import (
	"fmt"
)

// SchemaNotFoundError represents an error that occurs when a specified schema
// is not found in the library.
type SchemaNotFoundError struct {
	SchemaName string
	Err        error
}

// Error conforms to go conventions.
func (e SchemaNotFoundError) Error() string {
	return fmt.Sprintf(
		"could not locate the following schema in the Library: %s",
		e.SchemaName,
	)
}

// Unwrap conforms to go conventions.
func (e SchemaNotFoundError) Unwrap() error {
	return e.Err
}

// DatabaseError represents an error that occurs during a database operation.
type DatabaseError struct {
	Err error
}

// Error conforms to go conventions.
func (e DatabaseError) Error() string {
	return fmt.Sprintf(
		"database error occurred: %s",
		e.Err,
	)
}

// Unwrap conforms to go conventions.
func (e DatabaseError) Unwrap() error {
	return e.Err
}
