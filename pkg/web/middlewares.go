package web

import (
	"net/http"
)

type loginManager struct {
	loginUrl string
}

func (lm *loginManager) wrap(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, found := store.Get(r)
		if !found {
			http.Redirect(w, r, lm.loginUrl, http.StatusFound)
			return
		}

		h(w, r)
	}
}
