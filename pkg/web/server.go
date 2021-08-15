package web

import (
	"github.com/go-chi/chi"
)

func NewRouter() *chi.Mux {
	router := chi.NewRouter()

	lm := loginManager{
		loginUrl: "/login",
	}

	// User interface
	router.Get("/", lm.wrap(indexHandler))
	router.Get("/login", loginHandler)
	router.Post("/login", loginHandler)
	router.Get("/logout", logoutHandler)

	// RESTful resources
	router.Route("/profiles", func(r chi.Router) {
		r.Get("/", lm.wrap(voidHandler))
		r.Post("/", lm.wrap(voidHandler))

		r.Get("/{id}", lm.wrap(voidHandler))
		r.Delete("/{id}", lm.wrap(voidHandler))

		r.Route("/devices", func(r chi.Router) {
			r.Get("/", lm.wrap(voidHandler))
			r.Post("/", lm.wrap(voidHandler))

			r.Get("/{id}", lm.wrap(voidHandler))
			r.Delete("/{id}", lm.wrap(voidHandler))
		})
	})

	return router
}
