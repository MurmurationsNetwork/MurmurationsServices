package service

import (
	"testing"

	"github.com/MurmurationsNetwork/MurmurationsServices/services/cronjob/parseschema/internal/adapter/redisadapter"
	"github.com/go-playground/assert/v2"
)

var svc = NewSchemaService(redisadapter.NewClient())

func TestGetSchemaURL(t *testing.T) {
	url := svc.GetSchemaURL("test1")
	assert.Equal(t, "/schemas/test1.json", url)
}
