package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TODO: Not able to do unit tests.
func TestGetLinkedSchemas(t *testing.T) {
	_, ok := getLinkedSchemas(nil)
	assert.Equal(t, false, ok)

	_, ok = getLinkedSchemas(map[string]interface{}{})
	assert.Equal(t, false, ok)

	_, ok = getLinkedSchemas(map[string]interface{}{"profileUrl": "https://ic3.dev/test2.json"})
	assert.Equal(t, false, ok)

	_, ok = getLinkedSchemas(map[string]interface{}{"linkedSchemas": "https://ic3.dev/test2.json"})
	assert.Equal(t, false, ok)

	_, ok = getLinkedSchemas(map[string]interface{}{"linkedSchemas": false})
	assert.Equal(t, false, ok)

	_, ok = getLinkedSchemas(map[string]interface{}{"linkedSchemas": []string{}})
	assert.Equal(t, false, ok)
}
