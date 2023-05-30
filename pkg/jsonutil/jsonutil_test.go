package jsonutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToJSON(t *testing.T) {
	assert.Equal(t, map[string]interface{}{}, ToJSON(""))

	assert.Equal(
		t,
		"Demo Schema",
		ToJSON("{\"title\": \"Demo Schema\",\"y\": 6}")["title"],
	)

	assert.Equal(t, 6.0, ToJSON("{\"title\": \"Demo Schema\",\"x\": 6}")["x"])
}
