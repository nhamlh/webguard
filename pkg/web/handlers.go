package web

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/nhamlh/wg-dash/pkg/db"
)

func voidHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Simply works")
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFS(staticAssets, "static/templates/index.tpl"))
	t.Execute(w, nil)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFS(staticAssets, "static/templates/login.tpl"))

	switch r.Method {
	case http.MethodGet:
		err := t.Execute(w, nil)
		if err != nil {
			w.Write([]byte(fmt.Sprintf("Error reading login.tpl", err.Error())))
		}
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
			t.Execute(w, templateData{Errors: []string{"Invalid email or password"}})
		}

		session, err := sessionStore.New()
		if err != nil {
			w.WriteHeader(500)
			t.Execute(w, templateData{Errors: []string{"Internal error", err.Error()}})
		}
		session.Value = email

		cookie, err := sessionStore.Marshal(session)
		if err != nil {
			w.WriteHeader(500)
			t.Execute(w, templateData{Errors: []string{"Internal error", err.Error()}})
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

		sessionStore.Print()

		http.Redirect(w, r, r.Referer(), 302)
	default:
		w.WriteHeader(405)
	}
}
