package jsonutil

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestToJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected map[string]interface{}
	}{
		{
			name:     "valid JSON",
			input:    `{"key":"value"}`,
			expected: map[string]interface{}{"key": "value"},
		},
		{
			name:     "empty JSON",
			input:    `{}`,
			expected: map[string]interface{}{},
		},
		{
			name:     "invalid JSON",
			input:    `{"key":"value`,
			expected: map[string]interface{}{},
		},
		{
			name:  "nested JSON",
			input: `{"key1":"value1", "key2":{"nestedKey":"nestedValue"}}`,
			expected: map[string]interface{}{
				"key1": "value1",
				"key2": map[string]interface{}{"nestedKey": "nestedValue"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToJSON(tt.input)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestToString(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]interface{}
		expected string
	}{
		{
			name:     "valid map",
			input:    map[string]interface{}{"key": "value"},
			expected: `{"key":"value"}`,
		},
		{
			name:     "empty map",
			input:    map[string]interface{}{},
			expected: `{}`,
		},
		{
			name: "nested map",
			input: map[string]interface{}{
				"key1": "value1",
				"key2": map[string]interface{}{"nestedKey": "nestedValue"},
			},
			expected: `{"key1":"value1","key2":{"nestedKey":"nestedValue"}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToString(tt.input)
			require.JSONEq(t, tt.expected, result)
		})
	}
}
