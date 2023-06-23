package index_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/index"
)

func TestDatabaseError(t *testing.T) {
	want := "Encountered a DB error. Test message: Test error"
	err := index.DatabaseError{
		Message: "Test message",
		Err:     errors.New("Test error"),
	}

	require.Equal(
		t, want, err.Error(),
		"DatabaseError.Error() does not match expected",
	)

	var dbErr index.DatabaseError
	require.True(
		t, errors.As(err, &dbErr),
		"Unable to unwrap to DatabaseError",
	)
}

func TestNotFoundError(t *testing.T) {
	want := "Record not found: Test error"
	err := index.NotFoundError{
		Err: errors.New("Test error"),
	}

	require.Equal(
		t, want, err.Error(),
		"NotFoundError.Error() does not match expected",
	)

	var nfErr index.NotFoundError
	require.True(
		t, errors.As(err, &nfErr),
		"Unable to unwrap to NotFoundError",
	)
}

func TestValidationError(t *testing.T) {
	want := "Validation failed on field 'Test field': Test reason"
	err := index.ValidationError{
		Field:  "Test field",
		Reason: "Test reason",
	}

	require.Equal(
		t, want, err.Error(),
		"ValidationError.Error() does not match expected",
	)
}

func TestDeleteNodeError(t *testing.T) {
	err := index.DeleteNodeError{
		Message:    "Node cannot be deleted",
		Detail:     "Node is associated with active transactions",
		ProfileURL: "https://example.com/profile",
		NodeID:     "12345",
	}

	expected :=
		"Message: Node cannot be deleted, " +
			"Detail: Node is associated with active transactions, " +
			"Profile URL: https://example.com/profile, Node ID: 12345"

	require.Equal(
		t, expected, err.Error(),
		"DeleteNodeError.Error() does not match expected",
	)
}
