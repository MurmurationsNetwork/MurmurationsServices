package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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

	_, ok = getLinkedSchemas(map[string]interface{}{"linked_schemas": []string{"https://ic3.dev/test2.json"}})
	assert.Equal(t, false, ok)

	_, ok = getLinkedSchemas(map[string]interface{}{"linked_schemas": []interface{}{"https://ic3.dev/test2.json"}})
	assert.Equal(t, true, ok)

	linkedSchemas, ok := getLinkedSchemas(map[string]interface{}{"linked_schemas": []interface{}{"https://ic3.dev/test2.json"}})
	assert.Equal(t, "https://ic3.dev/test2.json", linkedSchemas[0])
}

func TestGetSchemaURL(t *testing.T) {
	url := getSchemaURL("test1.json")
	assert.Equal(t, "/schemas/test1.json.json", url)
}
