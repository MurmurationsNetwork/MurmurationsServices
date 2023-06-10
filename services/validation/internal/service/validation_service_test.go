package service

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetLinkedSchemas(t *testing.T) {
	tests := []struct {
		name          string
		data          interface{}
		expectedOk    bool
		expectedValue []string
	}{
		{
			name:       "Test with nil data",
			data:       nil,
			expectedOk: false,
		},
		{
			name:       "Test with empty map",
			data:       map[string]interface{}{},
			expectedOk: false,
		},
		{
			name: "Test with wrong key",
			data: map[string]interface{}{
				"profile_url": "https://ic3.dev/test2.json",
			},
			expectedOk: false,
		},
		{
			name: "Test with wrong type of value",
			data: map[string]interface{}{
				"linked_schemas": "https://ic3.dev/test2.json",
			},
			expectedOk: false,
		},
		{
			name: "Test with false value",
			data: map[string]interface{}{
				"linked_schemas": false,
			},
			expectedOk: false,
		},
		{
			name: "Test with empty slice",
			data: map[string]interface{}{
				"linked_schemas": []string{},
			},
			expectedOk: false,
		},
		{
			name: "Test with string slice",
			data: map[string]interface{}{
				"linked_schemas": []string{"https://ic3.dev/test2.json"},
			},
			expectedOk: false,
		},
		{
			name: "Test with interface slice",
			data: map[string]interface{}{
				"linked_schemas": []interface{}{"https://ic3.dev/test2.json"},
			},
			expectedOk:    true,
			expectedValue: []string{"https://ic3.dev/test2.json"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			linkedSchemas, ok := getLinkedSchemas(tt.data)
			require.Equal(t, tt.expectedOk, ok)
			if tt.expectedOk {
				require.Equal(t, tt.expectedValue, linkedSchemas)
			}
		})
	}
}
