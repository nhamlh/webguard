package web

import (
	"fmt"
	"github.com/go-chi/chi"
	"net/http"
	"time"
)

func StartServer() {
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

	srv := &http.Server{
		Handler: router,
		Addr:    "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Println("Starting server")
	srv.ListenAndServe()
}
