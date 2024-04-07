package httpclient

import "net/http"

// HttpClient defines the interface for an HTTP client
type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}
