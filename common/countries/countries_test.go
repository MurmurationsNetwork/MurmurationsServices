package countries

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFindAlpha2ByName(t *testing.T) {
	// Create a mock HTTP server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Return a sample response
		response := `{"TD": ["chad"]}`
		_, err := w.Write([]byte(response))
		require.NoError(t, err)
	}))
	defer mockServer.Close()

	url := mockServer.URL

	var country interface{} = "Chad"
	countryIso, err := FindAlpha2ByName(url, country)
	require.Equal(t, "TD", countryIso)
	require.Equal(t, nil, err)
}
