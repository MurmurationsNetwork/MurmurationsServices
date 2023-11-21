package service

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetLinkedSchemas(t *testing.T) {
	tests := []struct {
		name          string
		data          string
		expectedError string
		expectedValue []string
	}{
		{
			name:          "Test with empty JSON object",
			data:          "{}",
			expectedError: "linked schemas not found in profile",
		},
		{
			name:          "Test with wrong key",
			data:          `{"profile_url": "https://ic3.dev/test2.json"}`,
			expectedError: "linked schemas not found in profile",
		},
		{
			name:          "Test with wrong type of value",
			data:          `{"linked_schemas": "https://ic3.dev/test2.json"}`,
			expectedError: "linked schemas is not an array",
		},
		{
			name:          "Test with non-array value",
			data:          `{"linked_schemas": false}`,
			expectedError: "linked schemas is not an array",
		},
		{
			name:          "Test with empty array",
			data:          `{"linked_schemas": []}`,
			expectedError: "empty linked schemas array",
		},
		{
			name:          "Test with valid data",
			data:          `{"linked_schemas": ["https://ic3.dev/test2.json"]}`,
			expectedValue: []string{"https://ic3.dev/test2.json"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			linkedSchemas, err := getLinkedSchemas(tt.data)

			if tt.expectedError != "" {
				require.Error(t, err)
				require.Equal(t, tt.expectedError, err.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedValue, linkedSchemas)
			}
		})
	}
}
