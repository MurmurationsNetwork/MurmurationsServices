package countries

import (
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestFindAlpha2ByName(t *testing.T) {
	var country interface{} = "Chad"
	url := "https://raw.githubusercontent.com/MurmurationsNetwork/MurmurationsLibrary/main/countries/map.json"
	countryIso, err := FindAlpha2ByName(url, country)
	assert.Equal(t, "TD", countryIso)
	assert.Equal(t, nil, err)
}
