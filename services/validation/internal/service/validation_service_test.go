package service

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetLinkedSchemas(t *testing.T) {
	_, ok := getLinkedSchemas(nil)
	require.Equal(t, false, ok)

	_, ok = getLinkedSchemas(map[string]interface{}{})
	require.Equal(t, false, ok)

	_, ok = getLinkedSchemas(
		map[string]interface{}{"profile_url": "https://ic3.dev/test2.json"},
	)
	require.Equal(t, false, ok)

	_, ok = getLinkedSchemas(
		map[string]interface{}{"linked_schemas": "https://ic3.dev/test2.json"},
	)
	require.Equal(t, false, ok)

	_, ok = getLinkedSchemas(map[string]interface{}{"linked_schemas": false})
	require.Equal(t, false, ok)

	_, ok = getLinkedSchemas(
		map[string]interface{}{"linked_schemas": []string{}},
	)
	require.Equal(t, false, ok)

	_, ok = getLinkedSchemas(
		map[string]interface{}{
			"linked_schemas": []string{"https://ic3.dev/test2.json"},
		},
	)
	require.Equal(t, false, ok)

	_, ok = getLinkedSchemas(
		map[string]interface{}{
			"linked_schemas": []interface{}{"https://ic3.dev/test2.json"},
		},
	)
	require.Equal(t, true, ok)

	linkedSchemas, ok := getLinkedSchemas(
		map[string]interface{}{
			"linked_schemas": []interface{}{"https://ic3.dev/test2.json"},
		},
	)
	require.Equal(t, true, ok)
	require.Equal(t, "https://ic3.dev/test2.json", linkedSchemas[0])
}
