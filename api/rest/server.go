package rest

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/tempcke/rpm/usecase"
)

// Server is used to expose application over a restful API
type Server struct {
	http.Handler
	propRepo usecase.PropertyRepository
}

// NewServer constructs a Server
func NewServer(propRepo usecase.PropertyRepository) *Server {
	server := new(Server)
	server.propRepo = propRepo
	server.initRouter()
	return server
}

func (s *Server) initRouter() {
	r := chi.NewRouter()
	r.Route("/property", func(r chi.Router) {
		r.Post("/", addProperty(s.propRepo))
		r.Get("/", listProperties(s.propRepo))
		r.Route("/{propertyID}", func(r chi.Router) {
			r.Get("/", getProperty(s.propRepo))
			r.Delete("/", deleteProperty(s.propRepo))
		})
	})
	s.Handler = r
}
