package service

import (
	"io"
	"net/http"

	"github.com/MurmurationsNetwork/MurmurationsServices/services/dataproxy/config"
)

// SchemasResponse holds the fetched and validated schemas.
type SchemasResponse struct {
	JSONSchemas []string
	SchemaNames []string
}

// parseSchemas fetches and returns the JSON content for a list of schema names,
// including a default schema.
func parseSchemas(schemas []string) (*SchemasResponse, error) {
	// Include default schema and append provided schemas.
	schemaNames := append([]string{"default-v2.0.0"}, schemas...)
	jsonSchemas := make([]string, len(schemaNames))
	baseURL := config.Conf.Library.InternalURL + "/v2/schemas"

	for i, schema := range schemaNames {
		resp, err := http.Get(baseURL + "/" + schema)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		jsonSchemas[i] = string(body)
	}

	return &SchemasResponse{
		JSONSchemas: jsonSchemas,
		SchemaNames: schemaNames,
	}, nil
}
