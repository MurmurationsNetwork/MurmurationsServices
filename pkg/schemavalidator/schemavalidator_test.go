package schemavalidator_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/schemavalidator"
)

var StrSchema = `
{
	"$schema": "https://json-schema.org/draft-07/schema#",
	"$id": "https://test-library.murmurations.network/v2/schemas/people_schema-v0.1.0",
	"title": "People Schema",
	"type": "object",
	"properties": {
		"linked_schemas": {
			"title": "Linked Schemas",
			"type": "array",
			"items": {
				"type": "string",
				"pattern": "^[a-z][a-z0-9_]{7,97}-v[0-9]+\\.[0-9]+\\.[0-9]+$"
			},
			"minItems": 1,
			"uniqueItems": true
		},
		"name": {
			"title": "Full Name",
			"type": "string"
		},
		"primary_url": {
			"title": "Primary URL",
			"type": "string",
			"maxLength": 2000,
			"pattern": "^https?://.*"
		},
		"tags": {
			"title": "Tags",
			"type": "array",
			"items": {
				"type": "string"
			},
			"uniqueItems": true
		},
		"description": {
			"title": "Description/Bio",
			"type": "string"
		},
		"image": {
			"title": "Photo/Avatar",
			"type": "string",
			"maxLength": 2000,
			"pattern": "^https?://.*"
		},
		"images": {
			"title": "Other Images",
			"type": "array",
			"items": {
				"type": "object",
				"properties": {
					"name": {
						"title": "Image Name",
						"description": "Description of the image",
						"type": "string",
						"minLength": 1,
						"maxLength": 100
					},
					"url": {
						"title": "URL",
						"description": "A URL of the image starting with http:// or https://",
						"type": "string",
						"maxLength": 2000,
						"pattern": "^https?://.*"
					}
				},
				"required": [
					"url"
				]
			}
		},
		"urls": {
			"title": "Website Addresses/URLs",
			"type": "array",
			"items": {
				"type": "object",
				"properties": {
					"name": {
						"title": "URL Name",
						"type": "string"
					},
					"url": {
						"title": "URL",
						"type": "string",
						"maxLength": 2000,
						"pattern": "^https?://.*"
					}
				},
				"required": [
					"url"
				]
			},
			"uniqueItems": true
		},
		"knows_language": {
			"title": "Languages Spoken",
			"type": "array",
			"items": {
				"type": "string"
			},
			"minItems": 1,
			"uniqueItems": true
		},
		"contact_details": {
			"title": "Contact Details",
			"type": "object",
			"properties": {
				"email": {
					"title": "Email Address",
					"type": "string"
				},
				"contact_form": {
					"title": "Contact Form",
					"type": "string",
					"pattern": "^https?://.*"
				}
			}
		},
		"telephone": {
			"title": "Telephone Number",
			"type": "string"
		},
		"street_address": {
			"title": "Street Address",
			"type": "string"
		},
		"locality": {
			"title": "Locality",
			"type": "string"
		},
		"region": {
			"title": "Region",
			"type": "string"
		},
		"postal_code": {
			"title": "Postal Code",
			"type": "string"
		},
		"country_name": {
			"title": "Country name",
			"type": "string"
		},
		"country_iso_3166": {
			"title": "Country (2 letters)",
			"type": "string",
			"enum": [
				"AD",
				"AE",
				"AF"
			]
		},
		"geolocation": {
			"title": "Geolocation Coordinates",
			"type": "object",
			"properties": {
				"lat": {
					"title": "Latitude",
					"type": "number",
					"minimum": -90,
					"maximum": 90
				},
				"lon": {
					"title": "Longitude",
					"type": "number",
					"minimum": -180,
					"maximum": 180
				}
			},
			"required": [
				"lat",
				"lon"
			]
		},
		"relationships": {
			"title": "Relationships",
			"type": "array",
			"items": {
				"type": "object",
				"properties": {
					"predicate": {
						"title": "Predicate",
						"type": "string"
					},
					"object_url": {
						"title": "Object URL",
						"type": "string",
						"maxLength": 2000,
						"pattern": "^https?://.*"
					}
				},
				"required": [
					"predicate",
					"object_url"
				]
			},
			"uniqueItems": true
		}
	}
}
`

func TestValidate(t *testing.T) {
	tests := []struct {
		name     string
		profile  string
		expected bool
	}{
		{
			name:     "Valid Empty Profile",
			profile:  `{}`,
			expected: true,
		},
		// {
		// 	name:     "Overly Long String for Name",
		// 	profile:  `{"name": "` + strings.Repeat("a", 5000) + `"}`,
		// 	expected: false,
		// },
		{
			name:     "Invalid URL Format",
			profile:  `{"primary_url": "ftp://example.com"}`,
			expected: false,
		},
		{
			name:     "Valid URL Format",
			profile:  `{"primary_url": "https://example.com"}`,
			expected: true,
		},
		{
			name:     "Invalid URL in Images Array",
			profile:  `{"images": [{"name": "Test", "url": "ftp://invalid.com"}]}`,
			expected: false,
		},
		{
			name:     "Valid Image URL but Missing Required Field",
			profile:  `{"images": [{"name": "Test"}]}`,
			expected: false,
		},
		{
			name:     "Valid Geolocation",
			profile:  `{"geolocation": {"lat": 45.123, "lon": -75.123}}`,
			expected: true,
		},
		{
			name:     "Invalid Geolocation (Latitude Out of Bounds)",
			profile:  `{"geolocation": {"lat": 95.123, "lon": -75.123}}`,
			expected: false,
		},
		{
			name:     "Invalid Linked Schemas Pattern",
			profile:  `{"linked_schemas": ["InvalidSchema-v0.1.0"]}`,
			expected: false,
		},
		{
			name:     "Valid Linked Schemas Pattern",
			profile:  `{"linked_schemas": ["validschema_12345-v0.1.0"]}`,
			expected: true,
		},
		{
			name:     "Invalid ISO 3166 Country Code",
			profile:  `{"country_iso_3166": "USA"}`,
			expected: false,
		},
		{
			name:     "Valid ISO 3166 Country Code",
			profile:  `{"country_iso_3166": "AE"}`,
			expected: true,
		},
		// {
		// 	name:     "Overly Long String for Description",
		// 	profile:  `{"description": "` + strings.Repeat("a", 5000) + `"}`,
		// 	expected: false,
		// },
		// {
		// 	name:     "Invalid Email in Contact Details",
		// 	profile:  `{"contact_details": {"email": "invalid.email", "contact_form": "https://valid.com"}}`,
		// 	expected: false,
		// },
		// {
		// 	name:     "Non-UTF8 Characters in Name",
		// 	profile:  `{"name": "` + "\xC3\x28" + `"}`, // An invalid UTF-8 sequence
		// 	expected: false,
		// },
		// {
		// 	name:     "SQL Injection Attempt in Name",
		// 	profile:  `{"name": "Robert'); DROP TABLE Students;--"}`, // A basic SQL injection payload
		// 	expected: false,
		// },
		// {
		// 	name:     "NoSQL Injection Attempt in Name",
		// 	profile:  `{"name": "{$ne: null}"}`, // A basic MongoDB NoSQL injection payload
		// 	expected: false,
		// },
		// {
		// 	name:     "NoSQL Injection Attempt in Description",
		// 	profile:  `{"description": "{$gt: ''}"}`, // Another MongoDB NoSQL injection payload
		// 	expected: false,
		// },
		// {
		// 	name:     "XSS Payload in Image Name",
		// 	profile:  `{"images": [{"name": "<script>alert('xss')</script>", "url": "https://valid.com/image.jpg"}]}`, // A basic XSS payload
		// 	expected: false,
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator, err := schemavalidator.NewBuilder().
				WithStrSchemas([]string{StrSchema}).
				WithStrProfile(tt.profile).
				Build()

			require.NoError(t, err)
			result := validator.Validate()
			require.Equal(t, tt.expected, result.Valid)
		})
	}
}

// func FuzzValidate(f *testing.F) {
// 	f.Fuzz(func(t *testing.T, randomText string) {
// 		validator, err := schemavalidator.NewBuilder().
// 			WithStrSchemas([]string{StrSchema}).
// 			WithStrProfile(`{"name": "` + randomText + `"}`).
// 			Build()

// 		t.Logf("randomText: %s \n", randomText)

// 		require.NoError(t, err)
// 		result := validator.Validate()
// 		require.Equal(t, false, result.Valid)
// 	})
// }
