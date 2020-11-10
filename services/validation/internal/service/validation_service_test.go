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

	_, ok = getLinkedSchemas(map[string]interface{}{"profile_url": "https://ic3.dev/test2.json"})
	assert.Equal(t, false, ok)

	_, ok = getLinkedSchemas(map[string]interface{}{"linked_schemas": "https://ic3.dev/test2.json"})
	assert.Equal(t, false, ok)

	_, ok = getLinkedSchemas(map[string]interface{}{"linked_schemas": false})
	assert.Equal(t, false, ok)

	_, ok = getLinkedSchemas(map[string]interface{}{"linked_schemas": []string{}})
	assert.Equal(t, false, ok)
}
