package handler_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/handler"
)

func TestDeprecation(t *testing.T) {
	tests := []struct {
		name    string
		service string
	}{
		{
			"Test with LibraryAPI",
			"Library",
		},
		{
			"Test with indexAPI",
			"Index",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.Default()
			router.GET(
				"/deprecation",
				handler.NewDeprecationHandler(tt.service),
			)

			req, err := http.NewRequest(http.MethodGet, "/deprecation", nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)

			assert.Equal(t, http.StatusGone, resp.Code)

			expected := fmt.Sprintf(
				"The v1 API has been deprecated. "+
					"Please use the v2 API instead: "+
					"https://app.swaggerhub.com/apis-docs/MurmurationsNetwork/%sAPI/2.0.0",
				tt.service,
			)
			assert.Contains(t, resp.Body.String(), expected)
		})
	}
}

func TestEnsureFirstUpper(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "All lowercase",
			in:   "testinput",
			want: "Testinput",
		},
		{
			name: "All uppercase",
			in:   "TESTINPUT",
			want: "Testinput",
		},
		{
			name: "Mixed case",
			in:   "tEstInPut",
			want: "Testinput",
		},
		{
			name: "First uppercase, rest lowercase",
			in:   "Testinput",
			want: "Testinput",
		},
		{
			name: "Unicode input",
			in:   "téstInput",
			want: "Téstinput",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := handler.EnsureFirstUpper(tt.in)
			if got != tt.want {
				t.Errorf("EnsureFirstUpper() = %v, want %v", got, tt.want)
			}
		})
	}
}
