package rest

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/tempcke/rpm/actions"
	"github.com/tempcke/rpm/internal"
	"github.com/tempcke/rpm/internal/config"
)

// Server is used to expose application over a restful API
type Server struct {
	http.Handler
	Conf    config.Config
	actions actions.Actions
}

// NewServer constructs a Server
func NewServer(acts actions.Actions) *Server {
	server := new(Server)
	server.actions = acts
	server.InitRouter()
	server.Conf = config.GetConfig()
	return server
}
func (s Server) WithConfig(conf config.Config) *Server {
	s.Conf = conf
	s.InitRouter() // when conf changes router needs to be reloaded
	return &s
}

func (s *Server) InitRouter() {
	r := chi.NewRouter()

	// no auth
	r.Group(func(r chi.Router) {
		r.Get("/health", s.okHandler)
		r.Get("/health/ready", s.okHandler)
		r.Get("/health/live", s.okHandler)
	})

	ph := propertyHandler{actions: s.actions}

	// with auth
	r.Group(func(r chi.Router) {
		r.Use(s.AuthMW)
		r.Route("/property", func(r chi.Router) {
			r.Post("/", ph.addProperty)
			r.Get("/", ph.listProperties)
			r.Route("/{propertyID}", func(r chi.Router) {
				r.Put("/", ph.storeProperty)
				r.Get("/", ph.getProperty)
				r.Delete("/", ph.deleteProperty)
			})
		})
	})
	s.Handler = r
}
func (s *Server) okHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (s *Server) AuthMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			key       = s.Conf.GetString(internal.EnvAPIKey)
			secret    = s.Conf.GetString(internal.EnvAPISecret)
			reqKey    = r.Header.Get(HeaderAPIKey)
			reqSecret = r.Header.Get(HeaderAPISecret)
		)

		if key != "" && reqKey != key {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if secret != "" && reqSecret != secret {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
