package web

import (
	"fmt"
	"net/http"

	"github.com/nhamlh/wg-dash/pkg/db"
)

func voidHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Simply works")
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate("index", nil, w)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		renderTemplate("login", templateData{Errors: []string{"Invalid email or password"}}, w)
	case http.MethodPost:
		// already logged
		requestSession, found := sessionStore.Get(*r)
		if found && ! requestSession.IsExpired() {
			http.Redirect(w, r, "/", 301)

		}

		email := r.FormValue("email")
		password := r.FormValue("password")

		var user db.User
		db.DB.Get(&user, "SELECT * FROM users WHERE email=$1", email)

		if user == (db.User{}) || password != user.Password.String {
			w.WriteHeader(403)
			renderTemplate("login", templateData{Errors: []string{"Invalid email or password"}}, w)
		}

		session, err := sessionStore.New()
		if err != nil {
			w.WriteHeader(500)
			renderTemplate("login", templateData{Errors: []string{"Internal error"}}, w)
		}
		session.Value = email

		cookie, err := sessionStore.Marshal(session)
		if err != nil {
			w.WriteHeader(500)
			renderTemplate("login", templateData{Errors: []string{"Internal error"}}, w)
		}

		http.SetCookie(w, &cookie)
		http.Redirect(w, r, "/", 301)
	default:
		w.WriteHeader(405)
	}
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		session, found := sessionStore.Get(*r)
		if found {
			sessionStore.Delete(session.Id)
		}

		// sessionStore.Print()

		http.Redirect(w, r, r.Referer(), 302)
	default:
		w.WriteHeader(405)
	}
}
