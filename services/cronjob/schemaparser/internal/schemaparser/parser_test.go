package schemaparser_test

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/schemaparser/internal/model"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/schemaparser/internal/schemaparser"
)

var (
	peopleSchema = map[string]interface{}{
		"$schema": "https://json-schema.org/draft-07/schema#",
		"$id":     "https://test-library.murmurations.network/v2/schemas/people_schema-v0.1.0",
		"title":   "People Schema",
		"description": "A schema to add individuals in the regenerative " +
			"economy to the Murmurations Index",
		"type": "object",
		"properties": map[string]interface{}{
			"name": map[string]interface{}{
				"$ref":        "../fields/name.json",
				"title":       "Full Name",
				"description": "The full name of the person",
			},
			"member_of": map[string]interface{}{
				"title": "Member of",
				"type":  "array",
				"items": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"name": map[string]interface{}{
							"title":       "Name",
							"description": "The name of the entity.",
							"type":        "string",
						},
						"url": map[string]interface{}{
							"title":   "URL",
							"type":    "string",
							"pattern": "^https?://.*",
						},
					},
				},
			},
		},
		"required": []string{
			"linked_schemas",
			"name",
			"primary_url",
			"tags",
		},
		"metadata": map[string]interface{}{
			"creator": map[string]string{
				"name": "Murmurations Network",
				"url":  "https://murmurations.network/",
			},
			"schema": map[string]interface{}{
				"name":    "people_schema-v0.1.0",
				"purpose": "To map people within the regenerative economy",
				"url":     "https://murmurations.network",
			},
		},
	}
	nameSchema = map[string]interface{}{
		"title": "Name",
		"description": "The name of the entity, organization, " +
			"project, item, etc.",
		"type": "string",
		"metadata": map[string]interface{}{
			"creator": map[string]string{
				"name": "Murmurations Network",
				"url":  "https://murmurations.network",
			},
			"field": map[string]interface{}{
				"name":    "name",
				"version": "1.0.0",
			},
			"context": []string{
				"https://schema.org/name",
			},
			"purpose": "The common name that is generally used to refer to " +
				"the entity, organization, project, item, etc., which can be a " +
				"living being, a legal entity, an object (real or virtual) or " +
				"even a good or service.",
		},
	}
	expectedSchema = &model.SchemaJSON{
		Title: "People Schema",
		Description: "A schema to add individuals in the regenerative " +
			"economy to the Murmurations Index",
		Metadata: model.Metadata{
			Schema: model.InnerSchema{
				Name:    "people_schema-v0.1.0",
				Version: 0,
				URL:     "https://murmurations.network",
			},
		},
	}
	expectedFullJSON = bson.D{
		{Key: "$schema", Value: "https://json-schema.org/draft-07/schema#"},
		{
			Key:   "$id",
			Value: "https://test-library.murmurations.network/v2/schemas/people_schema-v0.1.0",
		},
		{Key: "title", Value: "People Schema"},
		{
			Key: "description",
			Value: "A schema to add individuals in the regenerative " +
				"economy to the Murmurations Index",
		},
		{Key: "type", Value: "object"},
		{Key: "properties", Value: bson.D{
			{Key: "name", Value: bson.D{
				{Key: "title", Value: "Full Name"},
				{Key: "description", Value: "The full name of the person"},
				{Key: "type", Value: "string"},
				{Key: "metadata", Value: bson.D{
					{Key: "creator", Value: bson.D{
						{Key: "name", Value: "Murmurations Network"},
						{Key: "url", Value: "https://murmurations.network"},
					}},
					{Key: "field", Value: bson.D{
						{Key: "name", Value: "name"},
						{Key: "version", Value: "1.0.0"},
					}},
					{Key: "context", Value: []interface{}{
						"https://schema.org/name",
					}},
					{
						Key: "purpose",
						Value: "The common name that is generally used to refer to the " +
							"entity, organization, project, item, etc., which can be " +
							"a living being, a legal entity, an object (real or virtual) " +
							"or even a good or service.",
					},
				}},
			}},
			{Key: "member_of", Value: bson.D{
				{Key: "title", Value: "Member of"},
				{Key: "type", Value: "array"},
				{Key: "items", Value: bson.D{
					{Key: "type", Value: "object"},
					{Key: "properties", Value: bson.D{
						{Key: "name", Value: bson.D{
							{Key: "title", Value: "Name"},
							{
								Key:   "description",
								Value: "The name of the entity.",
							},
							{Key: "type", Value: "string"},
						}},
						{Key: "url", Value: bson.D{
							{Key: "title", Value: "URL"},
							{Key: "type", Value: "string"},
							{Key: "pattern", Value: "^https?://.*"},
						}},
					}},
				}},
			}},
		}},
		{
			Key: "required",
			Value: []interface{}{
				"linked_schemas",
				"name",
				"primary_url",
				"tags",
			},
		},
		{Key: "metadata", Value: bson.D{
			{Key: "creator", Value: bson.D{
				{Key: "name", Value: "Murmurations Network"},
				{Key: "url", Value: "https://murmurations.network/"},
			}},
			{Key: "schema", Value: bson.D{
				{Key: "name", Value: "people_schema-v0.1.0"},
				{
					Key:   "purpose",
					Value: "To map people within the regenerative economy",
				},
				{Key: "url", Value: "https://murmurations.network"},
			}},
		}},
	}
)

func TestGetSchema(t *testing.T) {
	// Create a test server with handler.
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			var (
				resp        map[string]interface{}
				schemaBytes []byte
				err         error
			)

			if r.URL.Path == "/schemas/people_schema-v0.1.0.json" {
				schemaBytes, err = json.Marshal(peopleSchema)
			} else if r.URL.Path == "/fields/name.json" {
				schemaBytes, err = json.Marshal(nameSchema)
			}
			require.NoError(t, err)

			resp = map[string]interface{}{
				"content": base64.StdEncoding.EncodeToString(schemaBytes),
			}
			w.Header().Set("Content-Type", "application/json")
			require.NoError(t, json.NewEncoder(w).Encode(resp))
		}))
	defer ts.Close()

	fieldListMap := map[string]string{
		"name.json": ts.URL + "/fields/name.json",
	}
	schemaParser := schemaparser.NewSchemaParser(fieldListMap)
	result, err := schemaParser.GetSchema(
		ts.URL + "/schemas/people_schema-v0.1.0.json",
	)

	require.NoError(t, err)
	require.Equal(t, expectedSchema, result.Schema)
	require.Equal(t, toMap(expectedFullJSON), toMap(result.FullJSON))
}

func toMap(d bson.D) map[string]interface{} {
	m := make(map[string]interface{})
	for _, pair := range d {
		key := pair.Key
		switch v := pair.Value.(type) {
		case bson.D:
			m[key] = toMap(v)
		default:
			m[key] = v
		}
	}
	return m
}
