package cryptoutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSHA256(t *testing.T) {
	assert.Equal(
		t,
		"b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9",
		GetSHA256("hello world"),
	)
	assert.Equal(
		t,
		"ab413dfeee7bdb1aa0f9b7fcef435cee5b183b022c5e719e575ca9e03d51b709",
		GetSHA256("murmurations"),
	)
}
