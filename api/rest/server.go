package rest

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/tempcke/rpm/usecase"
)

type Server struct {
	http.Handler
	propRepo usecase.PropertyRepository
}

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
		r.Route("/{propertyID}", func(r chi.Router) {
			r.Get("/", getProperty(s.propRepo))
			// 	r.Put("/", putProperty)
			r.Delete("/", deleteProperty(s.propRepo))
		})
	})
	s.Handler = r
}
