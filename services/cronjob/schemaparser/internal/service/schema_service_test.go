package service

import (
	"testing"

	"github.com/go-playground/assert/v2"
)

var svc = schemaService{}

func TestGetSchemaURL(t *testing.T) {
	url := svc.getSchemaURL("test1")
	assert.Equal(t, "/schemas/test1.json", url)
}
