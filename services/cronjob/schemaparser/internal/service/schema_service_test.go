package service

import (
	"testing"
	"time"

	"github.com/go-playground/assert/v2"
)

var svc = schemaService{}

func TestGetSchemaURL(t *testing.T) {
	url := svc.getSchemaURL("test1")
	assert.Equal(t, "/schemas/test1.json", url)
}

func TestSetLastCommit(t *testing.T) {
	oldLastCommitTime := "2021-02-19T00:00:00Z"
	newLastCommitTime := "2021-02-19T00:00:00Z"

	t1, err := time.Parse(time.RFC3339, oldLastCommitTime)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	t2, err := time.Parse(time.RFC3339, newLastCommitTime)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	assert.Equal(t, true, int(t2.Sub(t1).Seconds()) < 180)
}

func TestSetLastCommit2(t *testing.T) {
	oldLastCommitTime := "2021-02-19T00:00:00Z"
	newLastCommitTime := "2021-02-19T00:04:00Z"

	t1, err := time.Parse(time.RFC3339, oldLastCommitTime)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	t2, err := time.Parse(time.RFC3339, newLastCommitTime)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	assert.Equal(t, false, int(t2.Sub(t1).Seconds()) < 180)
}
