package middlewares

import (
	"net/http"
	"time"

	"github.com/CSPF-Founder/api-scanner/code/panel/logger"
)

func LoggingMiddleware(appLogger *logger.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Serve the request to the next handler
			next.ServeHTTP(w, r)

			// Calculate the duration of the request
			duration := time.Since(start)

			appLogger.Provider.Info().
				Str("method", r.Method).
				Str("url", r.URL.String()).
				Str("remote_addr", r.RemoteAddr).
				Str("user_agent", r.UserAgent()).
				Dur("duration", duration).
				Msg("HTTP request")

		})
	}
}
