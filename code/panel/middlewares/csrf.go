package middlewares

import (
	"net/http"
	"strings"

	"github.com/CSPF-Founder/api-scanner/code/panel/internal/sessions"
)

// CSRFMiddleware checks for a valid CSRF token on each POST request, except for excluded paths.
func CSRFMiddleware(session *sessions.SessionManager, errorHandler http.Handler, excludedPaths []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get the CSRF token
			csrfToken := session.GetCSRF(r.Context())
			if csrfToken == "" {
				session.GenerateCSRF(r.Context())
			}

			// Only check CSRF for POST/PUT/PATCH/DELETE requests
			if !(r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodPatch || r.Method == http.MethodDelete) {
				next.ServeHTTP(w, r)
				return
			}

			// Determine if the request is multipart
			if strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
				// Parse multipart form
				if err := r.ParseMultipartForm(10 << 20); err != nil { // 10 MB
					http.Error(w, "Error parsing multipart form", http.StatusBadRequest)
					return
				}
			} else {
				// Parse regular form
				if err := r.ParseForm(); err != nil {
					http.Error(w, "Error parsing form", http.StatusBadRequest)
					return
				}
			}

			// Skip CSRF check for excluded paths
			for _, path := range excludedPaths {
				if strings.HasPrefix(r.URL.Path, path) {
					next.ServeHTTP(w, r)
					return
				}
			}

			var inputToken string
			// if it is delete method, check csrf token in header
			if r.Method == http.MethodDelete {
				inputToken = r.Header.Get("X-CSRF-Token")
			} else {
				// Get the CSRF token from the form
				inputToken = r.FormValue("csrf_token")
			}

			if !session.ValidateCSRF(r.Context(), inputToken) {
				// Handle CSRF check failure
				errorHandler.ServeHTTP(w, r)
				return
			}

			// Call the next handler
			next.ServeHTTP(w, r)
		})
	}
}
