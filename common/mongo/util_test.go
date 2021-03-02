package mongo

import (
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestgetURI(t *testing.T) {
	expect := "mongodb://admin:password@localhost:27017"
	actual := GetURI("admin", "password", "localhost:27017")
	assert.Equal(t, actual, expect)
}
