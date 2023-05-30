package mongo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetURI(t *testing.T) {
	expect := "mongodb://admin:password@localhost:27017"
	actual := GetURI("admin", "password", "localhost:27017")
	assert.Equal(t, actual, expect)
}
