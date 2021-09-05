package web

import (
	"github.com/go-chi/chi"
	"github.com/nhamlh/wg-dash/pkg/sso"
	"github.com/nhamlh/wg-dash/pkg/wg"
)

func NewRouter(wgInt *wg.Device, p *sso.Oauth2Provider) *chi.Mux {
	router := chi.NewRouter()

	lm := loginManager{
		loginUrl: "/login",
	}

	h := NewHandlers(wgInt, p)

	router.Get("/", lm.wrap(h.Index))

	// Working with devices
	router.Get("/new_device", lm.wrap(h.DeviceAdd))
	router.Post("/new_device", lm.wrap(h.DeviceAdd))
	router.Get("/devices/{id}/download", lm.wrap(h.DeviceDownload))
	router.Get("/devices/{id}/delete", lm.wrap(h.DeviceDelete))

	// Session management
	router.Get("/login", h.Login)
	router.Post("/login", h.Login)
	router.Get("/login/oauth", h.OauthLogin)
	router.Get("/login/oauth/callback", h.OauthCallback)
	router.Get("/logout", h.Logout)

	return router
}
