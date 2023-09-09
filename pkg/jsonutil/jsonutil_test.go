package jsonutil_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/jsonutil"
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
			result := jsonutil.ToJSON(tt.input)
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
			result := jsonutil.ToString(tt.input)
			require.JSONEq(t, tt.expected, result)
		})
	}
}

func TestHash(t *testing.T) {
	tests := []struct {
		name          string
		originalJSON  any
		alternateJSON any
		expErr        bool
	}{
		{
			name:          "Basic JSON",
			originalJSON:  `{"a": 1, "b": 2}`,
			alternateJSON: `{"b": 2, "a": 1}`,
			expErr:        false,
		},
		{
			name: "Complex Nested JSON",
			originalJSON: `{
				"user": {
				  "id": 12345,
				  "name": "John Doe",
				  "email": "johndoe@example.com",
				  "addresses": [
					{
					  "street": "123 Main St",
					  "city": "Anytown",
					  "state": "CA",
					  "zip": "12345"
					},
					{
					  "street": "456 Elm St",
					  "city": "Othertown",
					  "state": "TX",
					  "zip": "67890"
					}
				  ],
				  "phoneNumbers": ["123-456-7890", "987-654-3210"],
				  "isActive": true,
				  "roles": ["admin", "user", "guest"],
				  "preferences": {
					"theme": "dark",
					"language": "en-US",
					"notifications": {
					  "email": true,
					  "sms": false,
					  "push": true
					}
				  }
				},
				"products": [
				  {
					"id": 1,
					"name": "Laptop",
					"brand": "Brand A",
					"price": 1000.50
				  },
				  {
					"id": 2,
					"name": "Phone",
					"brand": "Brand B",
					"price": 500.25
				  }
				],
				"orderIDs": [101, 102, 103, 104],
				"timestamps": [1627890123, 1627890456, 1627890789],
				"messages": [
				  "Welcome to our platform!",
				  "Your order has been shipped.",
				  "Thank you for your purchase."
				]
			  }`,
			alternateJSON: `{
				"products": [
				  {
					"id": 1,
					"name": "Laptop",
					"price": 1000.50,
					"brand": "Brand A"
				  },
				  {
					"id": 2,
					"brand": "Brand B",
					"name": "Phone",
					"price": 500.25
				  }
				],
				"messages": [
				  "Welcome to our platform!",
				  "Your order has been shipped.",
				  "Thank you for your purchase."
				],
				"orderIDs": [101, 102, 103, 104],
				"timestamps": [1627890123, 1627890456, 1627890789],
				"user": {
				  "name": "John Doe",
				  "email": "johndoe@example.com",
				  "phoneNumbers": ["123-456-7890", "987-654-3210"],
				  "id": 12345,
				  "roles": ["admin", "user", "guest"],
				  "addresses": [
					{
					  "street": "123 Main St",
					  "city": "Anytown",
					  "state": "CA",
					  "zip": "12345"
					},
					{
					  "street": "456 Elm St",
					  "city": "Othertown",
					  "state": "TX",
					  "zip": "67890"
					}
				  ],
				  "isActive": true,
				  "preferences": {
					"theme": "dark",
					"language": "en-US",
					"notifications": {
					  "email": true,
					  "sms": false,
					  "push": true
					}
				  }
				}
			  }`,
			expErr: false,
		},
		{
			name: "Map Format JSON",
			originalJSON: map[string]interface{}{
				"name":     "Alice",
				"age":      30,
				"location": "New York",
			},
			alternateJSON: map[string]interface{}{
				"location": "New York",
				"age":      30,
				"name":     "Alice",
			},
			expErr: false,
		},
		{
			name:          "Invalid JSON",
			originalJSON:  `{"a": 1, "b": 2`,
			alternateJSON: `{"a": 1, "b": 2`,
			expErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash1, err1 := jsonutil.Hash(tt.originalJSON)
			hash2, err2 := jsonutil.Hash(tt.alternateJSON)

			if tt.expErr {
				require.Error(t, err1)
				require.Error(t, err2)
			} else {
				require.NoError(t, err1)
				require.NoError(t, err2)
				require.Equal(t, hash1, hash2, "Hashes should be equal for %s and %s", tt.originalJSON, tt.alternateJSON)
			}
		})
	}
}
