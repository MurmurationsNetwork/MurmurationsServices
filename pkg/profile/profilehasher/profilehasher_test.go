package profilehasher_test

import (
	"net/http"
	"net/http/httptest"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/profile/profilehasher"
)

func createProfileServer(t *testing.T, response string) *httptest.Server {
	return httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := w.Write([]byte(response))
			require.NoError(t, err, "Failed to write response")
		}),
	)
}

func createLibraryServer(
	t *testing.T,
	responses map[string]string,
) *httptest.Server {
	return httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !strings.Contains(r.URL.Path, "/v2/schemas/") {
				w.WriteHeader(http.StatusBadRequest)
				_, err := w.Write([]byte("Invalid schema path"))
				require.NoError(t, err, "Failed to write response")
				return
			}
			schemaName := path.Base(r.URL.Path)
			if response, ok := responses[schemaName]; ok {
				_, err := w.Write([]byte(response))
				require.NoError(t, err, "Failed to write response")
			} else {
				w.WriteHeader(http.StatusNotFound)
				_, err := w.Write([]byte("Schema not found"))
				require.NoError(t, err, "Failed to write response")
			}
		}),
	)
}

func TestHash(t *testing.T) {
	tests := []struct {
		name            string
		profileResponse string
		libraryMappings map[string]string // Map of URL path to response
		expectedErr     error
		expectedHash    string
	}{
		{
			name: "Profile with Only Name",
			profileResponse: `{
				"name": "Test Organization"
			}`,
			libraryMappings: map[string]string{},
			expectedErr:     nil,
			expectedHash:    "46ad5935291b46fa2a052d965a9208af03c8a8611e5898b3a7c7d0e39649471b",
		},
		{
			name: "Profile with Redundant Field",
			profileResponse: `{
				"redundant_field": "Some value"
			}`,
			libraryMappings: map[string]string{},
			expectedErr:     nil,
			expectedHash:    "44136fa355b3678a1146ad16f7e8649e94fb4fc21fe77e8310c060f61caaff8a",
		},
		{
			name: "Profile with One Linked Schema",
			profileResponse: `{
				"linked_schemas": ["schema_one-v1.0.0"],
				"name": "Test Organization",
				"primary_url": "https://test-organization.com",
				"tags": ["Test Tag"]
			}`,
			libraryMappings: map[string]string{
				"schema_one-v1.0.0": `{
					"properties": {
						"name": {
							"type": "string"
						},
						"primary_url": {
							"type": "string"
						},
						"tags": {
							"type": "array"
						}
					}
				}`,
			},
			expectedErr:  nil,
			expectedHash: "9be6e1a378d770d4213c0a3df56cbb5fa7a2605ebbbceeaec2b41835ad23cb1f",
		},
		{
			name: "Profile with Multiple Linked Schemas",
			profileResponse: `{
				"linked_schemas": ["schema_one-v1.0.0", "schema_two-v1.0.0"],
				"name": "Multi-Schema Organization",
				"primary_url": "https://multi-schema-organization.com",
				"tags": ["Test Tag"],
				"additional_field": "Some additional data",
				"redundant_field": "Some value"
			}`,
			libraryMappings: map[string]string{
				"schema_one-v1.0.0": `{
					"properties": {
						"name": {
							"type": "string"
						},
						"primary_url": {
							"type": "string"
						},
						"tags": {
							"type": "array"
						}
					}
				}`,
				"schema_two-v1.0.0": `{
					"properties": {
						"additional_field": {
							"type": "string"
						}
					}
				}`,
			},
			expectedErr:  nil,
			expectedHash: "61f723f41546908316bd647b25be8aafaf1422981f66ad5e9f308d0687326fdc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			profileServer := createProfileServer(t, tt.profileResponse)
			defer profileServer.Close()

			libraryServer := createLibraryServer(t, tt.libraryMappings)
			defer libraryServer.Close()

			ph := profilehasher.New(profileServer.URL, libraryServer.URL)
			hash, err := ph.Hash()

			if tt.expectedErr != nil {
				require.Error(t, err)
				require.Equal(t, tt.expectedErr, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedHash, hash)
			}
		})
	}
}
