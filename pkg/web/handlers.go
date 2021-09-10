package web

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/go-chi/chi"
	"github.com/nhamlh/webguard/pkg/db"
	"github.com/nhamlh/webguard/pkg/sso"
	"github.com/nhamlh/webguard/pkg/wg"
	"golang.org/x/crypto/bcrypt"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

type Handlers struct {
	wg *wg.Device
	op *sso.Oauth2Provider
}

func NewHandlers(wgInt *wg.Device, sp *sso.Oauth2Provider) Handlers {
	return Handlers{wg: wgInt, op: sp}
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

	var devStatus []map[string]string
	for _, dev := range devices {
		id := dev.Id
		name := dev.Name
		peer, _ := h.wg.GetPeer(dev.PrivateKey.PublicKey())
		lastSeen := humanize.Time(peer.LastHandshakeTime)

		devStatus = append(devStatus, map[string]string{
			"id":       strconv.Itoa(id),
			"name":     name,
			"pubkey":   peer.PublicKey.String(),
			"lastSeen": lastSeen,
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
	// already logged
	requestSession, found := sessionStore.Get(*r)
	if found && !requestSession.IsExpired() {
		http.Redirect(w, r, "/", http.StatusFound)
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
			w.WriteHeader(http.StatusUnauthorized)
			renderTemplate("login", templateData{"errors": []string{"Invalid email or password"}}, w)
			return
		}

		session, err := sessionStore.New()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			renderTemplate("login", templateData{"errors": []string{"Error happened [101]"}}, w)
			return
		}
		session.Value = email

		cookie, err := sessionStore.Marshal(session)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			renderTemplate("login", templateData{"errors": []string{"Error happened [102]"}}, w)
			return
		}

		http.SetCookie(w, &cookie)
		http.Redirect(w, r, "/", http.StatusFound)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *Handlers) OauthLogin(w http.ResponseWriter, r *http.Request) {
	// already logged
	requestSession, found := sessionStore.Get(*r)
	if found && !requestSession.IsExpired() {
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
	// already logged
	requestSession, found := sessionStore.Get(*r)
	if found && !requestSession.IsExpired() {
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
		db.DB.Get(&user, "SELECT * FROM users WHERE email=$1", email)

		// Insert user into database if nonexist
		if user == (db.User{}) {
			_, err = db.DB.Exec(`
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
			db.DB.Get(&user, "SELECT * FROM users WHERE email=$1", email)
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
		session, err := sessionStore.New()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			renderTemplate("login", templateData{"errors": []string{"Internal error"}}, w)
			return
		}
		session.Value = email

		cookie, err := sessionStore.Marshal(session)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			renderTemplate("login", templateData{"errors": []string{"Error while creating cookie"}}, w)
			return
		}

		http.SetCookie(w, &cookie)
		http.Redirect(w, r, "/", http.StatusFound)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *Handlers) Logout(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// Clear session cookie
		cookie, err := sessionStore.Marshal(Session{})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			renderTemplate("login", templateData{"errors": []string{"Internal error"}}, w)
			return
		}

		http.SetCookie(w, &cookie)
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
		session, _ := sessionStore.Get(*r)

		var user db.User
		db.DB.Get(&user, "SELECT * FROM users WHERE email=$1", session.Value)

		name := r.FormValue("name")
		if name == "" {
			w.WriteHeader(http.StatusBadRequest)
			renderTemplate("device", templateData{
				"user":   user,
				"errors": []string{"Device name cannot be empty"}}, w)
			return
		}

		prikey, err := wgtypes.GeneratePrivateKey()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			renderTemplate("device", templateData{
				"user":   user,
				"errors": []string{err.Error()}}, w)
			return
		}

		var devices []db.Device
		db.DB.Select(&devices, `SELECT * FROM devices`)
		deviceNum := getAvailNum(devices)

		var allowedIps []string
		for _, pr := range h.wg.PeerRoutes {
			allowedIps = append(allowedIps, pr.String())
		}

		_, err = db.DB.Exec(`
INSERT INTO
devices(user_id, name, private_key, num, allowed_ips)
values ($1,$2,$3,$4,$5)
`, user.Id, name, prikey.String(), deviceNum, strings.Join(allowedIps, ","))

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			renderTemplate("device", templateData{
				"user":   user,
				"errors": []string{"Cannot create device [E101]"}}, w)
			return
		}

		var device db.Device
		db.DB.Get(&device, `SELECT * FROM devices where private_key=$1`, prikey.String())

		peerIp, _ := h.wg.AllocateIP(deviceNum)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			renderTemplate("device", templateData{
				"user":   user,
				"errors": []string{"Cannot create device [E102]"}}, w)
			return
		}

		peer, err := generatePeerConfig(device, peerIp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			renderTemplate("device", templateData{
				"user":   user,
				"errors": []string{"Cannot create device [E103]"}}, w)
			return
		}

		h.wg.AddPeer(peer)

		http.Redirect(w, r, "/", http.StatusFound)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *Handlers) DeviceDelete(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		id := chi.URLParam(r, "id")

		session, _ := sessionStore.Get(*r)
		var user db.User
		db.DB.Get(&user, "SELECT * FROM users WHERE email=$1", session.Value)

		var device db.Device
		db.DB.Get(&device, "SELECT * FROM devices WHERE id=$1 AND user_id=$2", id, user.Id)

		if device == (db.Device{}) {
			w.WriteHeader(http.StatusBadRequest)
			renderTemplate("error", templateData{
				"user":   user,
				"errors": []string{"Cannot delete such device"}}, w)
			return
		}

		_, found := h.wg.GetPeer(device.PrivateKey.PublicKey())
		if !found {
			w.WriteHeader(http.StatusBadRequest)
			renderTemplate("error", templateData{
				"user": user,
				"errors": []string{
					"Cannot find peer from interface",
					device.PrivateKey.PublicKey().String(),
				},
			}, w)
			return
		}

		removed := h.wg.RemovePeer(device.PrivateKey.PublicKey())
		if !removed {
			w.WriteHeader(http.StatusBadRequest)
			renderTemplate("error", templateData{
				"user": user,
				"errors": []string{
					"Cannot remove peer from interface",
					device.PrivateKey.PublicKey().String(),
				},
			}, w)
			return
		}

		db.DB.Exec("DELETE FROM devices WHERE id=$1", id)

		http.Redirect(w, r, "/", http.StatusFound)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *Handlers) DeviceDownload(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		session, _ := sessionStore.Get(*r)
		var user db.User
		db.DB.Get(&user, "SELECT * FROM users WHERE email=$1", session.Value)

		var device db.Device
		db.DB.Get(&device, "SELECT * FROM devices WHERE id=$1 AND user_id=$2", chi.URLParam(r, "id"), user.Id)

		if device == (db.Device{}) {
			w.WriteHeader(http.StatusBadRequest)
			renderTemplate("error", templateData{
				"user":   user,
				"errors": []string{"Cannot delete such device"}}, w)
			return
		}

		clientCfg := generateClientConfig(h.wg, device)

		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Content-Disposition", "attachment; filename=wg.conf")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(clientCfg))
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
