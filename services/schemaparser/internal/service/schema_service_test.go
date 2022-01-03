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

func TestShouldSetLastCommitTime(t *testing.T) {
	t.Run("oldLastCommitTime is empty", func(t *testing.T) {
		oldLastCommitTime := ""
		newLastCommitTime := "2021-02-19T00:04:00Z"
		ok, _ := shouldSetLastCommitTime(oldLastCommitTime, newLastCommitTime)
		assert.Equal(t, true, ok)
	})
	t.Run("newLastCommitTime is empty", func(t *testing.T) {
		oldLastCommitTime := "2021-02-19T00:00:00Z"
		newLastCommitTime := ""
		ok, _ := shouldSetLastCommitTime(oldLastCommitTime, newLastCommitTime)
		assert.Equal(t, false, ok)
	})
	t.Run("no time difference", func(t *testing.T) {
		oldLastCommitTime := "2021-02-19T00:00:00Z"
		newLastCommitTime := "2021-02-19T00:00:00Z"
		ok, _ := shouldSetLastCommitTime(oldLastCommitTime, newLastCommitTime)
		assert.Equal(t, false, ok)
	})
	t.Run("should not set last commit time", func(t *testing.T) {
		oldLastCommitTime := "2021-02-19T00:00:00Z"
		newLastCommitTime := "2021-02-19T00:05:00Z"
		ok, _ := shouldSetLastCommitTime(oldLastCommitTime, newLastCommitTime)
		assert.Equal(t, false, ok)
	})
	t.Run("should set last commit time", func(t *testing.T) {
		oldLastCommitTime := "2021-02-19T00:00:00Z"
		newLastCommitTime := "2021-02-19T00:11:00Z"
		ok, _ := shouldSetLastCommitTime(oldLastCommitTime, newLastCommitTime)
		assert.Equal(t, true, ok)
	})
}
