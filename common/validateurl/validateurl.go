package validateurl

import (
	"net/url"
	"strings"
)

func Validate(rawUrl string) (string, error) {
	if !strings.Contains(rawUrl, "http") {
		rawUrl = "https://" + rawUrl
	}
	u, err := url.ParseRequestURI(rawUrl)
	if err != nil {
		return "", err
	}

	// url manipulation
	// if the first four character is "www.", remove it.
	var host string
	if u.Host[:4] == "www." {
		host = u.Host[4:]
	} else {
		host = u.Host
	}

	validatedUrl := host

	// url path
	if u.Path != "" {
		validatedUrl += u.Path
		if validatedUrl[(len(validatedUrl)-1):] == "/" {
			validatedUrl = validatedUrl[:(len(validatedUrl) - 1)]
		}
	}

	// url query
	if u.RawQuery != "" {
		validatedUrl += "?" + u.RawQuery
	}

	// if we have the "://", remove the right side of it
	position := strings.Index(validatedUrl, "://")
	if position != -1 {
		validatedUrl = validatedUrl[:position]
	}

	return validatedUrl, nil
}
