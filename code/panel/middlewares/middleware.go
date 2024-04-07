package middlewares

import (
	"encoding/json"
	"fmt"
	"net/http"

	ctx "github.com/CSPF-Founder/api-scanner/code/panel/context"
	"github.com/CSPF-Founder/api-scanner/code/panel/internal/sessions"
	"github.com/CSPF-Founder/api-scanner/code/panel/models"
)

// Use allows us to stack middleware to process the request
func Use(handler http.HandlerFunc, mid ...func(http.Handler) http.HandlerFunc) http.HandlerFunc {
	for _, m := range mid {
		handler = m(handler)
	}
	return handler
}

// ApplySecurityHeaders applies various security headers according to best-
// practices.
func ApplySecurityHeaders(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		csp := "frame-ancestors 'none';"
		w.Header().Set("Content-Security-Policy", csp)
		w.Header().Set("X-Frame-Options", "DENY")
		next.ServeHTTP(w, r)
	}
}

func GetContext(sessionManager *sessions.SessionManager) func(http.Handler) http.HandlerFunc {
	return func(handler http.Handler) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if id, ok := sessionManager.Get(r.Context(), "userID").(int64); ok {
				u, err := models.GetUserByID(id)
				if err != nil {
					r = ctx.Set(r, "user", nil)
				} else {
					r = ctx.Set(r, "user", u)
				}
			} else {
				r = ctx.Set(r, "user", nil)
			}
			handler.ServeHTTP(w, r)
		}
	}
}

// RequireLogin checks to see if the user is currently logged in.
// If not, the function returns a 302 redirect to the login page.
func RequireLogin(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if u := ctx.Get(r, "user"); u != nil {
			// If a password change is required for the user, then redirect them
			// to the login page
			handler.ServeHTTP(w, r)
			return
		}
		q := r.URL.Query()
		q.Set("next", r.URL.Path)
		http.Redirect(w, r, fmt.Sprintf("/users/login?%s", q.Encode()), http.StatusTemporaryRedirect)
	})
}

// EnforceViewOnly is a global middleware that limits the ability to edit
// objects to accounts with the PermissionModifyObjects permission.
func EnforceViewOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// If the request is for any non-GET HTTP method, e.g. POST, PUT,
		// or DELETE, we need to ensure the user has the appropriate
		// permission.
		if r.Method != http.MethodGet && r.Method != http.MethodHead && r.Method != http.MethodOptions {
			user := ctx.Get(r, "user").(models.User)
			access, err := user.HasPermission(models.PermissionModifyObjects)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			if !access {
				http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

// RequirePermission checks to see if the user has the requested permission
// before executing the handler. If the request is unauthorized, a JSONError
// is returned.
func RequirePermission(perm string) func(http.Handler) http.HandlerFunc {
	return func(next http.Handler) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			user := ctx.Get(r, "user").(models.User)
			access, err := user.HasPermission(perm)
			if err != nil {
				JSONError(w, http.StatusInternalServerError, err.Error())
				return
			}
			if !access {
				JSONError(w, http.StatusForbidden, http.StatusText(http.StatusForbidden))
				return
			}
			next.ServeHTTP(w, r)
		}
	}
}

// JSONError returns an error in JSON format with the given
// status code and message
func JSONError(w http.ResponseWriter, c int, m string) {
	cj, _ := json.MarshalIndent(models.Response{Success: false, Message: m}, "", "  ")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(c)
	fmt.Fprintf(w, "%s", cj)
}
