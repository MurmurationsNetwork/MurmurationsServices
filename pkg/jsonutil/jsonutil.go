package jsonutil

import (
	"encoding/json"
	"fmt"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/cryptoutil"
)

// ToJSON converts a JSON string to a map.
func ToJSON(s string) map[string]any {
	var raw map[string]any
	err := json.Unmarshal([]byte(s), &raw)
	if err != nil {
		return map[string]any{}
	}
	return raw
}

// ToString converts a map to a JSON string.
func ToString(m any) string {
	jsonString, err := json.Marshal(m)
	if err != nil {
		return ""
	}
	return string(jsonString)
}

// Hash computes the SHA-256 hash of the provided JSON data.
//
// The function ensures that different representations of the same JSON data
// (e.g., with fields in different orders) will produce the same hash.
// This is achieved by parsing the input JSON into a normalized format before hashing.
func Hash(data any) (string, error) {
	var doc []byte
	var err error

	switch v := data.(type) {
	case string:
		doc = []byte(v)
	case map[string]interface{}:
		doc, err = json.Marshal(v)
		if err != nil {
			return "", err
		}
	default:
		return "", fmt.Errorf("unsupported json data type")
	}

	// Normalize the JSON data.
	var parsedData interface{}
	if err := json.Unmarshal(doc, &parsedData); err != nil {
		return "", err
	}
	normalizedJSON, err := json.Marshal(parsedData)
	if err != nil {
		return "", err
	}

	return cryptoutil.ComputeSHA256(string(normalizedJSON)), nil
}
