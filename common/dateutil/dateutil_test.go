package dateutil

import (
	"testing"
	"time"

	"github.com/go-playground/assert/v2"
)

func TestGetNowUnix(t *testing.T) {
	assert.Equal(t, GetNowUnix(), time.Now().UTC().Unix())
}

func TestGetZeroValueUnix(t *testing.T) {
	assert.Equal(t, GetZeroValueUnix(), int64(-62135596800))
}

func TestNowSubtract(t *testing.T) {
	assert.Equal(t, NowSubtract(10*time.Second), time.Now().UTC().Add(-10*time.Second).Unix())
	assert.Equal(t, NowSubtract(600*time.Second), time.Now().UTC().Add(-10*time.Minute).Unix())
	assert.Equal(t, NowSubtract(86400*time.Second), time.Now().AddDate(0, 0, -1).UTC().Unix())
}
