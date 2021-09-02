package web

import (
	"fmt"
	"net/http"

	"strings"

	"strconv"

	"github.com/go-chi/chi"
	"github.com/nhamlh/wg-dash/pkg/db"
	"github.com/nhamlh/wg-dash/pkg/wg"
	"golang.org/x/crypto/bcrypt"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
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

	var devStatus []map[string]string
	for _, dev := range devices {
		id := dev.Id
		name := dev.Name
		prikey, _ := wgtypes.ParseKey(dev.PrivateKey)
		peer, _ := h.wg.GetPeer(prikey.PublicKey())
		lastSeen := peer.LastHandshakeTime.String()

		devStatus = append(devStatus, map[string]string{
			"id":       strconv.Itoa(id),
			"name":     name,
			"pubkey":   peer.PublicKey.String(),
			"lastSeen": lastSeen,
		})
	}

	data := templateData{
		"devices": devStatus,
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

func (h *Handlers) Device(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		renderTemplate("device", nil, w)
	case http.MethodPost:
		session, _ := sessionStore.Get(*r)

		name := r.FormValue("name")
		if name == "" {
			w.WriteHeader(http.StatusBadRequest)
			renderTemplate("device", templateData{"errors": []string{"Device name cannot be empty"}}, w)
			return
		}

		var user db.User
		db.DB.Get(&user, "SELECT * FROM users WHERE email=$1", session.Value)

		prikey, err := wgtypes.GeneratePrivateKey()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			renderTemplate("device", templateData{"errors": []string{err.Error()}}, w)
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
			renderTemplate("device", templateData{"errors": []string{"cannot insert device into database", err.Error()}}, w)
			return
		}

		var device db.Device
		db.DB.Get(&device, `SELECT * FROM devices where private_key=$1`, prikey.String())

		peerIp, _ := h.wg.AllocateIP(deviceNum)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			renderTemplate("device", templateData{"errors": []string{"cannot allocate IP for your device", err.Error()}}, w)
			return
		}

		peer, err := generatePeerConfig(device, peerIp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			renderTemplate("device", templateData{"errors": []string{"cannot import your device into server", err.Error()}}, w)
			return
		}

		h.wg.AddPeer(peer)

		http.Redirect(w, r, "/", http.StatusFound)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *Handlers) ClientConfig(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:

		var device db.Device
		db.DB.Get(&device, "SELECT * FROM devices WHERE id=$1", chi.URLParam(r, "id"))

		clientCfg := generateClientConfig(h.wg, device)

		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Content-Disposition", "attachment; filename=wg.conf")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(clientCfg))
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
