package service_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/MurmurationsNetwork/MurmurationsServices/services/validation/internal/service"
)

func TestNormalizeURL(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
		err      error
	}{
		{
			name:     "Regular URL",
			input:    "https://ic3.dev",
			expected: "ic3.dev",
		},
		{
			name:     "URL with trailing slash",
			input:    "https://www.ic3.dev/",
			expected: "ic3.dev",
		},
		{
			name:     "URL with www. at start",
			input:    "https://www.ic3.dev",
			expected: "ic3.dev",
		},
		{
			name:     "URL with www. in the middle",
			input:    "https://site.www.ic3.dev",
			expected: "site.www.ic3.dev",
		},
		{
			name:     "URL with query",
			input:    "https://www.ic3.dev/some/path/and/file.asp?id=123",
			expected: "ic3.dev/some/path/and/file.asp?id=123",
		},
		{
			name:     "URL with fragment",
			input:    "https://www.ic3.dev/page.html#section",
			expected: "ic3.dev/page.html#section",
		},
		{
			name:     "URL without protocol",
			input:    "ic3.dev/page.html",
			expected: "ic3.dev/page.html",
		},
		{
			name:     "URL without a valid top-level domain",
			input:    "https://max",
			expected: "",
			err:      errors.New("invalid URL"),
		},
		{
			name:     "URL with port number",
			input:    "https://ic3.dev:8000",
			expected: "ic3.dev:8000",
		},
		{
			name:     "URL with multiple levels of subdomains",
			input:    "https://sub.sub2.ic3.dev",
			expected: "sub.sub2.ic3.dev",
		},
		{
			name:     "URL with multiple paths and query parameters",
			input:    "https://www.ic3.dev/path/to/resource?id=123&order=desc",
			expected: "ic3.dev/path/to/resource?id=123&order=desc",
		},
		{
			name:     "URL with http protocol",
			input:    "http://ic3.dev",
			expected: "ic3.dev",
		},
		{
			name:     "Empty URL",
			input:    "",
			expected: "",
			err:      errors.New("invalid URL"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			normalizedURL, err := service.NormalizeURL(tc.input)
			require.Equal(t, tc.expected, normalizedURL)
			if tc.err != nil {
				require.Error(t, err)
				require.Equal(t, tc.err.Error(), err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
