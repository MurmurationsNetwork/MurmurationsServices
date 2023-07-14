package validateurl

import (
	"net/url"
	"strings"
)

func Validate(rawURL string) (string, error) {
	if !strings.Contains(rawURL, "http") {
		rawURL = "https://" + rawURL
	}
	u, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return "", err
	}

	// url manipulation
	// if the first four character is "www.", remove it.
	var host string

	if len(u.Host) >= 4 && u.Host[:4] == "www." {
		host = u.Host[4:]
	} else {
		host = u.Host
	}

	validatedURL := host

	// url path
	if u.Path != "" {
		validatedURL += u.Path
		if validatedURL[(len(validatedURL)-1):] == "/" {
			validatedURL = validatedURL[:(len(validatedURL) - 1)]
		}
	}

	// url query
	if u.RawQuery != "" {
		validatedURL += "?" + u.RawQuery
	}

	// if we have the "://", remove the right side of it
	position := strings.Index(validatedURL, "://")
	if position != -1 {
		validatedURL = validatedURL[:position]
	}

	return validatedURL, nil
}
