package library

import (
	"fmt"
)

// SchemaNotFoundError represents an error that occurs when a specified schema
// is not found in the library.
type SchemaNotFoundError struct {
	SchemaName string
}

// Error conforms to go conventions.
func (e SchemaNotFoundError) Error() string {
	return fmt.Sprintf(
		"could not locate the following schema in the Library: %s",
		e.SchemaName,
	)
}

// DatabaseError represents an error that occurs during a database operation.
type DatabaseError struct {
	Operation string
}

// Error conforms to go conventions.
func (e DatabaseError) Error() string {
	return fmt.Sprintf(
		"database error occurred during %s operation.",
		e.Operation,
	)
}
