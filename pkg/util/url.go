package util

import "net/url"

func IsValidUrl(u1 string) bool {
	_, err := url.ParseRequestURI(u1)
	if err != nil {
		return false
	}

	u, err := url.Parse(u1)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}

	// Check if the URL has a valid scheme (http or https)
	if u.Scheme != "http" && u.Scheme != "https" {
		return false
	}

	return true
}
