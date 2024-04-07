package context

import (
	"net/http"

	"context"
)

// Get retrieves a value from the request context
func Get(r *http.Request, key interface{}) interface{} {
	return r.Context().Value(key)
}

// Set stores a value on the request context
func Set(r *http.Request, key, val interface{}) *http.Request {
	if val == nil {
		return r
	}

	return r.WithContext(context.WithValue(r.Context(), key, val))
}
