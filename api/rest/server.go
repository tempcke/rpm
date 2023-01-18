package rest

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/tempcke/rpm/internal"
	"github.com/tempcke/rpm/internal/config"
	"github.com/tempcke/rpm/usecase"
)

// Server is used to expose application over a restful API
type Server struct {
	http.Handler
	propRepo usecase.PropertyRepository
	Conf     config.Config
}

// NewServer constructs a Server
func NewServer(propRepo usecase.PropertyRepository) *Server {
	server := new(Server)
	server.propRepo = propRepo
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
	r.Use(s.AuthMW)
	r.Route("/property", func(r chi.Router) {
		r.Post("/", addProperty(s.propRepo))
		r.Get("/", listProperties(s.propRepo))
		r.Route("/{propertyID}", func(r chi.Router) {
			r.Put("/", storeProperty(s.propRepo))
			r.Get("/", getProperty(s.propRepo))
			r.Delete("/", deleteProperty(s.propRepo))
		})
	})
	s.Handler = r
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
