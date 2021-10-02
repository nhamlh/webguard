package web

import (
	"net/http"
)

// RequireLoginAt is a middleware which require user to be logged in otherwise
// it redirects user to loginUrl
func RequireLoginAt(loginUrl string) func(next http.Handler) http.Handler {
	middleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, found := store.Get(r)
			if !found {
				http.Redirect(w, r, loginUrl, http.StatusFound)
				return
			}

			next.ServeHTTP(w, r)
		})
	}

	return middleware
}
