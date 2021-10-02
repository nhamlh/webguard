package web

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"
	"github.com/nhamlh/webguard/pkg/sso"
	wireguard "github.com/nhamlh/webguard/pkg/wg"
	"net/http"
)

func NewRouter(db *sqlx.DB, wgInt *wireguard.Interface, p *sso.Oauth2Provider) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	h := NewHandlers(db, wgInt, p)

	// session management
	r.Group(func(r chi.Router) {
		r.Get("/login", h.Login)
		r.Post("/login", h.Login)
		r.Get("/login/oauth", h.OauthLogin)
		r.Get("/login/oauth/callback", h.OauthCallback)
		r.Get("/logout", h.Logout)
	})

	// require login
	r.Group(func(r chi.Router) {
		r.Use(RequireLoginAt("/login"))

		r.Get("/", h.Index)

		r.Route("/devices", func(r chi.Router) {
			r.Get("/", h.DeviceAdd)
			r.Post("/", h.DeviceAdd)
			r.Get("/{id}/install", h.DeviceInstall)
			r.Get("/{id}/delete", h.DeviceDelete)
		})
	})

	// error pages
	r.Group(func(r chi.Router) {
		r.NotFound(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			renderTemplate("error", templateData{
				"errors": []string{"Route does not exist"}}, w)
		})
		r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusMethodNotAllowed)
			renderTemplate("error", templateData{
				"errors": []string{"Method is not valid"}}, w)
		})
	})

	return r
}
