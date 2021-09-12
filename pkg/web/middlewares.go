package web

import (
	"net/http"
	"time"

	"fmt"
	"github.com/nhamlh/webguard/pkg/db"
)

type loginManager struct {
	loginUrl string
}

func (lm *loginManager)wrap(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter,r *http.Request) {
		session, found := sessionStore.Get(*r)
		if !found {
			http.Redirect(w, r, lm.loginUrl, 301)
			return
		}

		if session.Expire.Before(time.Now()) {
			fmt.Println("Session expired")
			http.Redirect(w, r, lm.loginUrl, 301)
			return
		}

		var user db.User
		err := db.DB.Get(&user, "SELECT * FROM users WHERE email=$1", session.Value)
		if err != nil || user == (db.User{}) {
			http.Redirect(w, r, lm.loginUrl, http.StatusFound)
			return
		}

		h(w, r)
	}
}
