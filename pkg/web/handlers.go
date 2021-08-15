package web

import (
	"fmt"
	"net/http"

	"github.com/nhamlh/wg-dash/pkg/db"
	"golang.org/x/crypto/bcrypt"
)

func voidHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Simply works")
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	// doesn't need further check. LoginManager already
	// check for validity of session
	session, _ := sessionStore.Get(*r)

	var user db.User
	db.DB.Get(&user, "SELECT * FROM users WHERE email=$1", session.Value)

	var devices []db.Device
	db.DB.Select(&devices, "SELECT * FROM devices WHERE user_id=$1", user.Id)

	data := templateData{
		"devices": devices,
	}
	renderTemplate("index", data, w)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	// already logged
	requestSession, found := sessionStore.Get(*r)
	if found && !requestSession.IsExpired() {
		http.Redirect(w, r, "/", 301)
	}

	switch r.Method {
	case http.MethodGet:
		renderTemplate("login", nil, w)
	case http.MethodPost:
		email := r.FormValue("email")
		password := r.FormValue("password")

		var user db.User
		db.DB.Get(&user, "SELECT * FROM users WHERE email=$1", email)

		err := bcrypt.CompareHashAndPassword([]byte(user.Password.String), []byte(password))

		samePassword := true
		if err != nil {
			samePassword = false
		}

		if user == (db.User{}) || !samePassword {
			w.WriteHeader(403)
			renderTemplate("login", templateData{"errors": []string{"Invalid email or password"}}, w)
			return
		}

		session, err := sessionStore.New()
		if err != nil {
			w.WriteHeader(500)
			renderTemplate("login", templateData{"errors": []string{"Internal error"}}, w)
			return
		}
		session.Value = email

		cookie, err := sessionStore.Marshal(session)
		if err != nil {
			w.WriteHeader(500)
			renderTemplate("login", templateData{"errors": []string{"Internal error"}}, w)
			return
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

		http.Redirect(w, r, r.Referer(), 302)
	default:
		w.WriteHeader(405)
	}
}
