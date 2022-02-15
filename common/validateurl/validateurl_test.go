package validateurl

import (
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestValidate(t *testing.T) {
	// test regular url
	var url string = "https://ic3.dev"
	parsedUrl, _ := Validate(url)
	assert.Equal(t, "ic3.dev", parsedUrl)

	// if the last character is a slash, remove it
	url = "https://www.ic3.dev/"
	parsedUrl, _ = Validate(url)
	assert.Equal(t, "ic3.dev", parsedUrl)

	// if the string includes ://, discard this substring and everything to the left of it
	url = "https://www.ic3.dev/path1://path2"
	parsedUrl, _ = Validate(url)
	assert.Equal(t, "ic3.dev/path1", parsedUrl)

	// if the string starts with www., remove those 4 characters
	url = "https://www.ic3.dev"
	parsedUrl, _ = Validate(url)
	assert.Equal(t, "ic3.dev", parsedUrl)

	// www. in the middle of url
	url = "https://site.www.ic3.dev"
	parsedUrl, _ = Validate(url)
	assert.Equal(t, "site.www.ic3.dev", parsedUrl)

	// query in url
	url = "https://www.ic3.dev/some/path/and/file.asp?id=123"
	parsedUrl, _ = Validate(url)
	assert.Equal(t, "ic3.dev/some/path/and/file.asp?id=123", parsedUrl)

	// fragments in url
	url = "https://www.ic3.dev/page.html#section"
	parsedUrl, _ = Validate(url)
	assert.Equal(t, "ic3.dev/page.html#section", parsedUrl)
}
