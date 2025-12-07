package service

import (
	"errors"
	"net/url"
	"strings"
)

var ErrInvalidURL = errors.New("invalid URL")

// NormalizeURL function takes a raw URL string as input, validates it, and
// returns a normalized URL string or an error.
func NormalizeURL(rawURL string) (string, error) {
	u, err := validateURL(rawURL)
	if err != nil {
		return "", err
	}

	normalizedURL := normalizeURL(u)

	// If the normalizedURL is empty after all the manipulations, return an error.
	if normalizedURL == "" {
		return "", ErrInvalidURL
	}

	// Return the normalized and manipulated URL.
	return normalizedURL, nil
}

// normalizeURL function takes a raw URL string as input and
// returns a normalized URL string.
func normalizeURL(u *url.URL) string {
	// Remove "www." from the start of the host part of the URL, if it exists.
	var host string
	if len(u.Host) >= 4 && u.Host[:4] == "www." {
		host = u.Host[4:]
	} else {
		host = u.Host
	}

	// Start building the normalizedURL with the host part.
	normalizedURL := host

	// If a path exists in the URL, append it to the validatedURL.
	// Also, remove any trailing "/" from the validatedURL.
	if u.Path != "" {
		normalizedURL += u.Path
		if normalizedURL[(len(normalizedURL)-1):] == "/" {
			normalizedURL = normalizedURL[:(len(normalizedURL) - 1)]
		}
	}

	// If a query exists in the URL, append it to the normalizedURL.
	if u.RawQuery != "" {
		normalizedURL += "?" + u.RawQuery
	}

	// If "://" exists in the normalizedURL, remove everything after it.
	position := strings.Index(normalizedURL, "://")
	if position != -1 {
		normalizedURL = normalizedURL[:position]
	}

	return normalizedURL
}

// validateURL function takes a raw URL string as input, validates it, and
// returns a parsed URL or an error.
func validateURL(rawURL string) (*url.URL, error) {
	// Check if the rawURL is empty or only contains "https://" or "www."
	// If so, return an error as these are not valid URLs.
	if rawURL == "" || rawURL == "https://" || rawURL == "www." {
		return nil, ErrInvalidURL
	}

	// Normalize the scheme to lowercase for case-insensitive checking.
	// Check if the URL starts with "http://" or "https://" (case-insensitive).
	rawURLLower := strings.ToLower(rawURL)
	if strings.HasPrefix(rawURLLower, "http://") {
		// Find the position of "://" in the original URL (case-insensitive)
		schemeEnd := strings.Index(rawURL, "://")
		if schemeEnd != -1 {
			// Normalize http:// to lowercase
			rawURL = "http://" + rawURL[schemeEnd+3:]
		}
	} else if strings.HasPrefix(rawURLLower, "https://") {
		// Find the position of "://" in the original URL (case-insensitive)
		schemeEnd := strings.Index(rawURL, "://")
		if schemeEnd != -1 {
			// Normalize https:// to lowercase
			rawURL = "https://" + rawURL[schemeEnd+3:]
		}
	} else if !strings.Contains(rawURLLower, "http") {
		// If the rawURL does not contain "http" (case-insensitive), prepend "https://" to it.
		rawURL = "https://" + rawURL
	}

	// Parse the rawURL to ensure it's a valid URL.
	// If parsing fails, return an error.
	u, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return nil, ErrInvalidURL
	}

	// Validate if the domain contains a period indicating a top-level domain.
	if !strings.Contains(u.Host, ".") {
		return nil, ErrInvalidURL
	}

	return u, nil
}
