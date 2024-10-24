package service

import (
	"fmt"
	"io"
	"net/http"

	"github.com/MurmurationsNetwork/MurmurationsServices/services/dataproxy/config"
)

// SchemasResponse holds the fetched and validated schemas.
type SchemasResponse struct {
	JSONSchemas []string
	SchemaNames []string
}

// DefaultSchema defines the default schema version to be fetched.
const DefaultSchema = "default-v2.1.0"

// ParseSchemas fetches and returns the JSON content for a list of schema names,
// including a default schema.
func ParseSchemas(schemas []string) (*SchemasResponse, error) {
	// Include default schema and append provided schemas.
	schemaNames := append([]string{DefaultSchema}, schemas...)
	jsonSchemas := make([]string, len(schemaNames))
	baseURL := fmt.Sprintf("%s/v2/schemas", config.Conf.Library.InternalURL)

	for i, schema := range schemaNames {
		schemaURL := fmt.Sprintf("%s/%s", baseURL, schema)
		resp, err := http.Get(schemaURL)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch schema '%s': %w", schema, err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read schema '%s': %w", schema, err)
		}
		jsonSchemas[i] = string(body)
	}

	return &SchemasResponse{
		JSONSchemas: jsonSchemas,
		SchemaNames: schemaNames,
	}, nil
}
