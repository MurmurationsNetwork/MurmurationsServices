package library_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/MurmurationsNetwork/MurmurationsServices/services/library/internal/library"
)

func TestSchemaNotFoundError(t *testing.T) {
	err := &library.SchemaNotFoundError{SchemaName: "test-schema"}

	expected := "could not locate the following schema in the Library: test-schema"
	require.Equal(
		t,
		expected,
		err.Error(),
		"SchemaNotFoundError Error() message was incorrect",
	)
}

func TestDatabaseError(t *testing.T) {
	err := &library.DatabaseError{Err: errors.New("db error")}

	expected := "database error occurred: db error"
	require.Equal(
		t,
		expected,
		err.Error(),
		"DatabaseError Error() message was incorrect",
	)
}
