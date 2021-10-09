package web

import (
	"net/http"

	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"
	"github.com/nhamlh/webguard/pkg/sso"
	wireguard "github.com/nhamlh/webguard/pkg/wg"
	"log"
	"time"
)

type Server struct {
	r  chi.Router
	db sqlx.DB
	wg wireguard.Interface
	op sso.Oauth2Provider
}

func NewServer(db sqlx.DB, wg wireguard.Interface, op sso.Oauth2Provider) Server {

	r := chi.NewRouter()
	s := Server{
		r:  r,
		db: db,
		wg: wg,
		op: op,
	}

	s.initRoutes()

	return s
}

func (s *Server) StartAt(host string, port int) {
	listen := fmt.Sprintf("%s:%d", host, port)

	srv := &http.Server{
		Handler:      s.r,
		Addr:         listen,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println("Web server is listening at", listen)
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(fmt.Errorf("Web server failed: %v", err))
	}
}

func (s *Server) initRoutes() {
	s.r.Use(middleware.Logger)

	h := NewHandlers(&s.db, &s.wg, &s.op)

	// session management
	s.r.Group(func(r chi.Router) {
		r.Get("/login", h.Login)
		r.Post("/login", h.Login)
		r.Get("/logout", h.Logout)

		if s.op != (sso.Oauth2Provider{}) {
			r.Get("/login/oauth", h.OauthLogin)
			r.Get("/login/oauth/callback", h.OauthCallback)
		}
	})

	// require login
	s.r.Group(func(r chi.Router) {
		r.Use(RequireLoginAt("/login"))

		r.Get("/", h.Index)

		r.Route("/devices", func(r chi.Router) {
			r.Get("/", h.DeviceAdd)
			r.Post("/", h.DeviceAdd)
			r.Get("/{id}/install", h.DeviceInstall)
			r.Get("/{id}/delete", h.DeviceDelete)
		})

		r.Get("/change_password", h.ChangePasswd)
		r.Post("/change_password", h.ChangePasswd)
	})

	// error pages
	s.r.Group(func(r chi.Router) {
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
}
