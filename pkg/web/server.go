package web

import (
	"github.com/go-chi/chi"
	"github.com/nhamlh/wg-dash/pkg/wg"
)

func NewRouterFor(wgInt *wg.Device) *chi.Mux {
	router := chi.NewRouter()

	lm := loginManager{
		loginUrl: "/login",
	}

	h := NewHandlers(wgInt)

	// User interface
	router.Get("/", lm.wrap(h.Index))
	router.Get("/new_device", lm.wrap(h.Device))
	router.Post("/new_device", lm.wrap(h.Device))
	router.Get("/devices/{id}/download", lm.wrap(h.ClientConfig))
	router.Get("/devices/{id}/delete", lm.wrap(h.DeleteDevice))
	router.Get("/login", h.Login)
	router.Post("/login", h.Login)
	router.Get("/logout", h.Logout)

	// RESTful resources
	router.Route("/profiles", func(r chi.Router) {
		r.Get("/", lm.wrap(h.Void))
		r.Post("/", lm.wrap(h.Void))

		r.Get("/{id}", lm.wrap(h.Void))
		r.Delete("/{id}", lm.wrap(h.Void))

		r.Route("/devices", func(r chi.Router) {
			r.Get("/", lm.wrap(h.Void))
			r.Post("/", lm.wrap(h.Device))

			r.Get("/{id}", lm.wrap(h.Void))
			r.Delete("/{id}", lm.wrap(h.Void))
		})
	})

	return router
}
