package countries

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestFindAlpha2ByName(t *testing.T) {
	// Create a mock HTTP server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Return a sample response
		response := `{"TD": ["chad"]}`
		w.Write([]byte(response))
	}))
	defer mockServer.Close()

	url := mockServer.URL

	var country interface{} = "Chad"
	countryIso, err := FindAlpha2ByName(url, country)
	assert.Equal(t, "TD", countryIso)
	assert.Equal(t, nil, err)
}
