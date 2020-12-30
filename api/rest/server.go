package rest

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi"
	log "github.com/sirupsen/logrus"
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
		r.Post("/", s.addProperty)
		r.Route("/{propertyID}", func(r chi.Router) {
			r.Get("/", s.getProperty)
			// 	r.Put("/", putProperty)
			// 	r.Delete("/", deleteProperty)
		})
	})
	s.Handler = r
}

func (s Server) addProperty(w http.ResponseWriter, r *http.Request) {
	data := PropertyModel{}
	if err := s.decodeRequestData(w, r.Body, &data); err != nil {
		return
	}

	uc := usecase.NewAddPropertyCommand(s.propRepo)
	property := s.propRepo.NewProperty(
		data.Street, data.City, data.State, data.Zip,
	)

	if err := uc.Execute(property); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Missing or invalid fields"))
		return
	}

	w.WriteHeader(http.StatusCreated)
	s.jsonResponse(w, NewPropertyModel(property))
}

func (s Server) getProperty(w http.ResponseWriter, r *http.Request) {
	propertyID := chi.URLParam(r, "propertyID")
	uc := usecase.NewGetPropertyQuery(s.propRepo)
	property, err := uc.Execute(propertyID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Invalid property ID"))
		return
	}
	s.jsonResponse(w, NewPropertyModel(property))
}

func (s Server) decodeRequestData(w http.ResponseWriter, body io.Reader, data interface{}) error {
	err := json.NewDecoder(body).Decode(&data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Request body was not valid json"))
		return err
	}
	return nil
}

func (s Server) jsonResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Println("json.Encode response failed: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}
}
