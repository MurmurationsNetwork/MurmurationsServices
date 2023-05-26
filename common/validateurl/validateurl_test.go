package validateurl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	// test regular url
	url := "https://ic3.dev"
	parsedURL, _ := Validate(url)
	assert.Equal(t, "ic3.dev", parsedURL)

	// if the last character is a slash, remove it
	url = "https://www.ic3.dev/"
	parsedURL, _ = Validate(url)
	assert.Equal(t, "ic3.dev", parsedURL)

	// if the string includes ://, discard this substring and everything to the left of it
	url = "https://www.ic3.dev/path1://path2"
	parsedURL, _ = Validate(url)
	assert.Equal(t, "ic3.dev/path1", parsedURL)

	// if the string starts with www., remove those 4 characters
	url = "https://www.ic3.dev"
	parsedURL, _ = Validate(url)
	assert.Equal(t, "ic3.dev", parsedURL)

	// www. in the middle of url
	url = "https://site.www.ic3.dev"
	parsedURL, _ = Validate(url)
	assert.Equal(t, "site.www.ic3.dev", parsedURL)

	// query in url
	url = "https://www.ic3.dev/some/path/and/file.asp?id=123"
	parsedURL, _ = Validate(url)
	assert.Equal(t, "ic3.dev/some/path/and/file.asp?id=123", parsedURL)

	// fragments in url
	url = "https://www.ic3.dev/page.html#section"
	parsedURL, _ = Validate(url)
	assert.Equal(t, "ic3.dev/page.html#section", parsedURL)

	// url without protocol
	url = "ic3.dev/page.html"
	parsedURL, _ = Validate(url)
	assert.Equal(t, "ic3.dev/page.html", parsedURL)
}
