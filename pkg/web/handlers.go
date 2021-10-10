package web

import (
	"fmt"
	"log"
	"net/http"

	// "github.com/dustin/go-humanize"
	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
	"github.com/nhamlh/webguard/pkg/db"
	"github.com/nhamlh/webguard/pkg/session"
	"github.com/nhamlh/webguard/pkg/sso"
	"github.com/nhamlh/webguard/pkg/wg"
)

var store = session.NewSessionStore()

type Handlers struct {
	db *sqlx.DB
	wg *wg.Interface
	op *sso.Oauth2Provider
}

func NewHandlers(db *sqlx.DB, wgInt *wg.Interface, sp *sso.Oauth2Provider) Handlers {
	return Handlers{db: db, wg: wgInt, op: sp}
}

func (h *Handlers) Void(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Simply works")
}

func (h *Handlers) Index(w http.ResponseWriter, r *http.Request) {
	s, _ := store.Get(r)
	user := s.Values["user"].(db.User)

	var devices []db.Device
	h.db.Select(&devices, "SELECT * FROM devices WHERE user_id=$1", user.Id)

	var devStatus []map[string]interface{}
	for _, dev := range devices {
		status, peer := dev.Status(*h.wg)

		devStatus = append(devStatus, map[string]interface{}{
			"dev":  dev,
			"stat": status.String(),
			"peer": peer,
		})
	}

	// Get help section. helps is a list because
	// a param might appears multiple times in the url
	helps := r.URL.Query()["help"]
	var help string
	if len(helps) > 0 {
		help = helps[0]
	}

	data := templateData{
		"user":    user,
		"devices": devStatus,
		"help":    help,
	}

	renderTemplate("index", data, w)
}

func (h *Handlers) Login(w http.ResponseWriter, r *http.Request) {
	// already logged in
	_, found := store.Get(r)
	if found {
		http.Redirect(w, r, "/", http.StatusFound)
	}

	switch r.Method {
	case http.MethodGet:
		renderTemplate("login", nil, w)
	case http.MethodPost:
		email := r.FormValue("email")
		password := r.FormValue("password")

		user, found := db.GetUserByEmail(email, *h.db)
		if !found || !user.PasswdMatched([]byte(password)) {
			w.WriteHeader(http.StatusUnauthorized)
			renderTemplate("login", templateData{"errors": []string{"Invalid email or password"}}, w)
			return
		}

		isFirstLogin := user.LastLogin.IsZero()
		user.RecordLogin()
		user.Save(*h.db)

		session := store.New()
		session.Values["user"] = user
		store.Save(*session, w)

		if isFirstLogin {
			http.Redirect(w, r, "/change_password", http.StatusFound)
			return
		}

		http.Redirect(w, r, "/", http.StatusFound)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *Handlers) OauthLogin(w http.ResponseWriter, r *http.Request) {
	// already logged in
	_, found := store.Get(r)
	if found {
		http.Redirect(w, r, "/", http.StatusFound)
	}

	switch r.Method {
	case http.MethodGet:
		h.op.Redirect(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *Handlers) OauthCallback(w http.ResponseWriter, r *http.Request) {
	// already logged in
	_, found := store.Get(r)
	if found {
		http.Redirect(w, r, "/", http.StatusFound)
	}

	switch r.Method {
	case http.MethodGet:
		token, err := h.op.GetToken(r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			renderTemplate("login", templateData{"errors": []string{
				"Cannot retrieve access token from authorization server",
			}}, w)
			return
		}

		// getting user email
		email := h.op.Email(*token)
		if email == "" {
			renderTemplate("login", templateData{"errors": []string{
				"Cannot retrieve your email from authorization server",
			}}, w)
			return
		}

		var user db.User
		h.db.Get(&user, "SELECT * FROM users WHERE email=$1", email)

		// Insert user into database if nonexist
		if user == (db.User{}) {
			_, err = h.db.Exec(`
insert into
users(email,password,is_admin,auth_type)
values($1, "", 0, 1)`, email)

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				renderTemplate("login", templateData{"errors": []string{
					"User creation failed",
					email,
					err.Error(),
				}}, w)
				return
			}

			// let's try again
			h.db.Get(&user, "SELECT * FROM users WHERE email=$1", email)
			if user == (db.User{}) {
				w.WriteHeader(http.StatusInternalServerError)
				renderTemplate("login", templateData{"errors": []string{
					"User creation failed",
					email,
				}}, w)
				return
			}
		}

		// All good, generate a session and push to client
		session := store.New()
		session.Values["user"] = user

		store.Save(*session, w)

		http.Redirect(w, r, "/", http.StatusFound)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *Handlers) Logout(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		store.Destroy(session.Session{}, w)
		http.Redirect(w, r, r.Referer(), http.StatusFound)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *Handlers) DeviceAdd(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		renderTemplate("device", nil, w)
	case http.MethodPost:
		session, _ := store.Get(r)
		user := session.Values["user"].(db.User)

		name := r.FormValue("name")
		if name == "" {
			w.WriteHeader(http.StatusBadRequest)
			renderTemplate("device", templateData{
				"user":   user,
				"errors": []string{"Device name cannot be empty"}}, w)
			return
		}

		var devices []db.Device
		h.db.Select(&devices, `SELECT * FROM devices`)
		devNum := getAvailNum(devices)

		dev, err := db.NewDevice(user.Id, name, devNum, h.wg.PeerRoutes)
		if err != nil {
			log.Println(fmt.Errorf("Cannot create device for user %d: %v", user.Id, err))
			w.WriteHeader(http.StatusInternalServerError)
			renderTemplate("device", templateData{
				"user":   user,
				"errors": []string{"Cannot create device [E101]"}}, w)
			return
		}

		if err := dev.Save(*h.db); err != nil {
			log.Println(fmt.Errorf("Cannot save device %s to db: %v", dev.PrivateKey.PublicKey(), err))
			w.WriteHeader(http.StatusInternalServerError)
			renderTemplate("device", templateData{
				"user":   user,
				"errors": []string{"Cannot create device [E102]"}}, w)
			return
		}

		if err := dev.AddTo(h.wg); err != nil {
			log.Println(fmt.Errorf("Cannot add device %s to wg interface: %v", dev.PrivateKey.PublicKey(), err))
			w.WriteHeader(http.StatusInternalServerError)
			renderTemplate("device", templateData{
				"user":   user,
				"errors": []string{"Device created but can't be activated [E103]"}}, w)
			return
		}

		http.Redirect(w, r, "/", http.StatusFound)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *Handlers) DeviceDelete(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		id := chi.URLParam(r, "id")

		session, _ := store.Get(r)
		user := session.Values["user"].(db.User)

		var device db.Device
		h.db.Get(&device, "SELECT * FROM devices WHERE id=$1 AND user_id=$2", id, user.Id)

		if device == (db.Device{}) {
			w.WriteHeader(http.StatusBadRequest)
			renderTemplate("notif", templateData{
				"user":   user,
				"errors": []string{"Cannot delete such device [101]"}}, w)
			return
		}

		if err := device.RemoveFrom(h.wg); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			renderTemplate("notif", templateData{
				"user":   user,
				"errors": []string{"Cannot delete such device [102]"}}, w)
			return
		}

		if _, err := h.db.Exec("DELETE FROM devices WHERE id=$1", id); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			renderTemplate("notif", templateData{
				"user":   user,
				"errors": []string{"Cannot delete such device [103]"}}, w)
			return
		}

		http.Redirect(w, r, "/", http.StatusFound)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// DeviceInstall renders installation page for a device, also generates qrcode
// and client configuration file
func (h *Handlers) DeviceInstall(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		session, _ := store.Get(r)
		user := session.Values["user"].(db.User)

		var device db.Device
		h.db.Get(&device, "SELECT * FROM devices WHERE id=$1 AND user_id=$2", chi.URLParam(r, "id"), user.Id)

		if device == (db.Device{}) {
			w.WriteHeader(http.StatusBadRequest)
			renderTemplate("notif", templateData{
				"user":   user,
				"errors": []string{"Cannot load such device"}}, w)
			return
		}

		downloads := r.URL.Query()["dl"]
		if len(downloads) > 0 {
			clientCfg := device.GenClientConfig(h.wg)

			w.Header().Set("Content-Type", "text/plain")
			w.Header().Set("Content-Disposition", "attachment; filename=wg.conf")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(clientCfg))
			return
		}

		w.WriteHeader(http.StatusOK)
		renderTemplate("device_install", templateData{
			"user":         user,
			"download_url": fmt.Sprintf("%s?dl=1", r.URL.Path),
			"qrcode":       device.GenQRCode(h.wg)}, w)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *Handlers) ChangePasswd(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r)
	u := session.Values["user"].(db.User)

	switch r.Method {
	case http.MethodGet:
		renderTemplate("change_password", templateData{
			"user": u}, w)
	case http.MethodPost:
		current := r.FormValue("current_password")
		new := r.FormValue("new_password")

		if current == new {
			w.WriteHeader(http.StatusBadRequest)
			renderTemplate("change_password", templateData{
				"user":   u,
				"errors": []string{"New password must differ from your current password"}}, w)
			return
		}

		var user db.User
		h.db.Get(&user, "SELECT * FROM users WHERE id=$1", u.Id)

		if m := user.PasswdMatched([]byte(current)); !m {
			w.WriteHeader(http.StatusBadRequest)
			renderTemplate("change_password", templateData{
				"user":   user,
				"errors": []string{"Current password does not match"}}, w)
			return
		}

		if err := user.NewPasswd(new); err != nil {
			log.Println(err.Error())

			w.WriteHeader(http.StatusBadRequest)
			renderTemplate("change_password", templateData{
				"user":   user,
				"errors": []string{"Cannot update your password"}}, w)
			return
		}

		if err := user.Save(*h.db); err != nil {
			log.Println(err.Error())

			w.WriteHeader(http.StatusBadRequest)
			renderTemplate("change_password", templateData{
				"user":   user,
				"errors": []string{"Cannot update your password"}}, w)
			return
		}

		w.WriteHeader(http.StatusOK)
		renderTemplate("notif", templateData{
			"user": user,
			"msg":  []string{"Your password has changed!"}}, w)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
