package web

import (
	"fmt"
	"net/http"

	"github.com/nhamlh/wg-dash/pkg/db"
	"github.com/nhamlh/wg-dash/pkg/wg"
	"golang.org/x/crypto/bcrypt"
)

type Handlers struct {
	wg *wg.Device
}

func NewHandlers(wgInt *wg.Device) Handlers {
	return Handlers{wg: wgInt}
}

func (h *Handlers) Void(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Simply works")
}

func (h *Handlers) Index(w http.ResponseWriter, r *http.Request) {
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

func (h *Handlers) Login(w http.ResponseWriter, r *http.Request) {
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

func (h *Handlers) Logout(w http.ResponseWriter, r *http.Request) {
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
