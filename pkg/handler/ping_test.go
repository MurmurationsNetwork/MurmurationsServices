package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/handler"
)

func TestPing(t *testing.T) {
	// Set up the Gin router.
	router := gin.Default()

	// Register the Ping endpoint.
	router.GET("/ping", handler.PingHandler)

	// Create a request to the Ping endpoint.
	req, err := http.NewRequest(http.MethodGet, "/ping", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Record the response.
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Check the status code and body.
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, "pong!", resp.Body.String())
}
