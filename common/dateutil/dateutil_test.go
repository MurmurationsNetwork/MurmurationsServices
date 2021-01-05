package dateutil

import (
	"testing"
	"time"

	"github.com/go-playground/assert/v2"
)

func TestGetNowUnix(t *testing.T) {
	expect := time.Now().UTC().Unix()
	actual := GetNowUnix()
	assert.Equal(t, actual, expect)
}

func TestGetZeroValueUnix(t *testing.T) {
	expect := int64(-62135596800)
	actual := GetZeroValueUnix()
	assert.Equal(t, actual, expect)
}

func TestNowSubtract(t *testing.T) {
	expect := time.Now().UTC().Add(-10 * time.Second).Unix()
	actual := NowSubtract(10 * time.Second)
	assert.Equal(t, actual, expect)

	expect = time.Now().UTC().Add(-10 * time.Minute).Unix()
	actual = NowSubtract(600 * time.Second)
	assert.Equal(t, actual, expect)

	expect = time.Now().AddDate(0, 0, -1).UTC().Unix()
	actual = NowSubtract(86400 * time.Second)
	assert.Equal(t, actual, expect)
}
