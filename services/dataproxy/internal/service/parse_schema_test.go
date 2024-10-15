package service_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/MurmurationsNetwork/MurmurationsServices/services/dataproxy/config"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/dataproxy/internal/service"
)

// TestParseSchemas ensures that schemas are correctly fetched and parsed
// for both URL-based schema references and loaded JSON schemas.
func TestParseSchemas(t *testing.T) {
	// Setup a mock server to simulate schema fetching.
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		// Simulate different responses for the default schema and custom schemas.
		if strings.Contains(r.URL.Path, "default-v2.1.0") {
			_, err = w.Write([]byte(`{"default": "schema"}`))
		} else {
			_, err = w.Write([]byte(`{"custom": "schema"}`))
		}
		// Check for errors when writing the response.
		if err != nil {
			http.Error(w, "Failed to write response", http.StatusInternalServerError)
		}
	}))
	defer mockServer.Close()

	// Mock the config to use the test server's URL.
	config.Conf.Library.InternalURL = mockServer.URL

	// Define test cases.
	tests := []struct {
		name    string
		schemas []string
		want    *service.SchemasResponse
		wantErr bool
	}{
		{
			name:    "Valid schemas",
			schemas: []string{"custom-schema-1", "custom-schema-2"},
			wantErr: false,
			want: &service.SchemasResponse{
				JSONSchemas: []string{`{"default": "schema"}`, `{"custom": "schema"}`, `{"custom": "schema"}`},
				SchemaNames: []string{"default-v2.1.0", "custom-schema-1", "custom-schema-2"},
			},
		},
		{
			name:    "Empty schema list",
			schemas: []string{},
			wantErr: false,
			want: &service.SchemasResponse{
				JSONSchemas: []string{`{"default": "schema"}`},
				SchemaNames: []string{"default-v2.1.0"},
			},
		},
	}

	// Run the test cases.
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.ParseSchemas(tt.schemas)

			// Check for errors first.
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)

				// Compare the expected and actual results.
				assert.Equal(t, tt.want.JSONSchemas, got.JSONSchemas)
				assert.Equal(t, tt.want.SchemaNames, got.SchemaNames)
			}
		})
	}
}
