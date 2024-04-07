package utils

import (
	"net/http"
	"net/url"
	"regexp"
)

// IsRelativeURL checks if a URL is relative or not
func IsRelativeURL(inputURL string) bool {
	if inputURL == "/" {
		return true
	}

	parsedURL, err := url.Parse(inputURL)
	if err != nil {
		return false
	}

	if parsedURL.Host == "" && regexp.MustCompile(`^/[A-Za-z]`).MatchString(parsedURL.Path) {
		return true
	}

	return false
}

// GetRelativePath returns the relative path of the previous page
func GetRelativePath(r *http.Request) (previousPage string) {
	// strip the domain from the previous page
	previousPageURL, err := url.Parse(r.Header.Get("Referer"))
	if err == nil && previousPageURL.Path != "" {
		previousPage = previousPageURL.Path
	} else {
		previousPage = "/"
	}

	return previousPage
}

// RedirectBack redirects to the previous page
func RedirectBack(w http.ResponseWriter, r *http.Request) {
	previousPage := GetRelativePath(r)

	// Redirect to the referring URL
	http.Redirect(w, r, previousPage, http.StatusFound)
}
