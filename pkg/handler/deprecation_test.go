package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/handler"
)

func TestDeprecation(t *testing.T) {
	// Set up the Gin router.
	router := gin.Default()

	// Register the Deprecation endpoint.
	router.GET("/deprecation", handler.DeprecationHandler)

	// Create a request to the Deprecation endpoint.
	req, err := http.NewRequest(http.MethodGet, "/deprecation", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Record the response.
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Check the status code.
	assert.Equal(t, http.StatusGone, resp.Code)

	// Check part of the response body.
	// This could be more thorough, depending on the expected format of the error message.
	assert.Contains(
		t,
		resp.Body.String(),
		"The v1 API has been deprecated. Please use the v2 API instead",
	)
}
