package rest_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/controller/rest"
)

func TestNodeCreateRequest_Validate(t *testing.T) {
	tests := []struct {
		name     string
		request  rest.NodeCreateRequest
		mockHTTP func() *httptest.Server
		hasError bool
	}{
		{
			name: "empty profile URL",
			request: rest.NodeCreateRequest{
				ID:     "testNode2",
				Status: "Validated",
			},
			hasError: true,
		},
		{
			name: "unsupported URL scheme",
			request: rest.NodeCreateRequest{
				ID:         "testNode4",
				ProfileURL: "ftp://test.com",
				Status:     "Validated",
			},
			hasError: true,
		},
		{
			name: "invalid hostname",
			request: rest.NodeCreateRequest{
				ID:         "testNode5",
				ProfileURL: "http://invalid-hostname",
				Status:     "Validated",
			},
			hasError: true,
		},
		{
			name: "valid request",
			request: rest.NodeCreateRequest{
				ID:         "testNode1",
				ProfileURL: "https://test.com",
				Status:     "Validated",
			},
			mockHTTP: func() *httptest.Server {
				return httptest.NewServer(
					http.HandlerFunc(
						func(w http.ResponseWriter, _ *http.Request) {
							w.WriteHeader(http.StatusOK)
						},
					),
				)
			},
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock HTTP server.
			if tt.mockHTTP != nil {
				server := tt.mockHTTP()
				defer server.Close()
				tt.request.ProfileURL = server.URL
			}

			err := tt.request.Validate()
			if tt.hasError {
				require.NotEmpty(t, err, "Expected non-empty error slice")
			} else {
				require.Empty(t, err, "Expected empty error slice")
			}
		})
	}
}

func TestIsValidURL(t *testing.T) {
	tests := []struct {
		name  string
		url   string
		valid bool
	}{
		{
			name:  "unsupported URL scheme",
			url:   "ftp://test.com",
			valid: false,
		},
		{
			name:  "invalid hostname",
			url:   "http://invalid-hostname",
			valid: false,
		},
		{
			name:  "the data proxy app",
			url:   "http://data-proxy-app:8080",
			valid: true,
		},
		{
			name:  "valid URL with dot in hostname",
			url:   "http://test.com",
			valid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u, err := url.Parse(tt.url)
			require.NoError(t, err)

			valid := rest.IsValidURL(u)
			if tt.valid {
				require.True(t, valid, "Expected URL to be valid")
			} else {
				require.False(t, valid, "Expected URL to be invalid")
			}
		})
	}
}
