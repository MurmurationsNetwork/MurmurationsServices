package countries

import (
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestFindAlpha2ByName(t *testing.T) {
	var country interface{} = "Chad"
	countryIso, err := FindAlpha2ByName(country)
	assert.Equal(t, "TD", countryIso)
	assert.Equal(t, nil, err)
}
