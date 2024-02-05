package rest

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/tempcke/rpm/actions"
	oapi "github.com/tempcke/rpm/api/rest/openapi"
	"github.com/tempcke/rpm/entity"
	"github.com/tempcke/rpm/internal"
	"github.com/tempcke/rpm/internal/config"
	"github.com/tempcke/rpm/internal/lib/log"
)

const (
	HeaderAPIKey    = "X-Api-Key"
	HeaderAPISecret = "x-Api-Secret"
)

var _ oapi.ServerInterface = (*Server)(nil)

type Header struct{ k, v string }

type Server struct {
	Conf    config.Config
	actions actions.Actions
}

func (s *Server) LeaseProperty(w http.ResponseWriter, r *http.Request) {
	// TODO implement me
	panic("implement me")
}
func (s *Server) GetLease(w http.ResponseWriter, r *http.Request, id string) {
	// TODO implement me
	panic("implement me")
}
func (s *Server) ListLeases(w http.ResponseWriter, r *http.Request) {
	// TODO implement me
	panic("implement me")
}

func (s *Server) AddTenant(w http.ResponseWriter, r *http.Request) {
	s.StoreTenant(w, r, entity.NewID())
}
func (s *Server) StoreTenant(w http.ResponseWriter, r *http.Request, id string) {
	var (
		ctx     = r.Context()
		resCode = http.StatusCreated
		data    oapi.StoreTenantReq
	)
	if err := decodeRequestData(w, r.Body, &data); err != nil {
		return
	}

	if cur, _ := s.actions.GetTenant(ctx, id); cur != nil && cur.GetID() == id {
		resCode = http.StatusOK
	}

	tenant := data.Tenant.ToTenant().WithID(id)

	if _, err := s.actions.StoreTenant(ctx, tenant); err != nil {
		errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	jsonResponse(w, resCode, oapi.NewGetTenantRes(tenant),
		Header{"Location", "/tenant/" + tenant.ID})
}
func (s *Server) GetTenant(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()
	tenant, err := s.actions.GetTenant(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, internal.ErrEntityNotFound):
			errorResponse(w, http.StatusNotFound, err.Error())
		case errors.Is(err, internal.ErrInternal):
			errorResponse(w, http.StatusInternalServerError, err.Error())
		default:
			errorResponse(w, http.StatusBadRequest, err.Error())
		}
		return
	}
	jsonResponse(w, http.StatusOK, oapi.NewGetTenantRes(*tenant))
}
func (s *Server) ListTenants(w http.ResponseWriter, r *http.Request) {
	var ctx = r.Context()
	list, err := s.actions.ListTenants(ctx)
	if err != nil {
		s.logError(err)
		errorResponse(w, http.StatusInternalServerError, "Error fetching list")
		return
	}
	jsonResponse(w, http.StatusOK, oapi.ToTenantList(list...))
}

func (s *Server) AddProperty(w http.ResponseWriter, r *http.Request) {
	s.StoreProperty(w, r, entity.NewID())
}
func (s *Server) StoreProperty(w http.ResponseWriter, r *http.Request, id string) {
	var (
		ctx     = r.Context()
		resCode = http.StatusCreated
		data    oapi.StorePropertyReq
	)
	if err := decodeRequestData(w, r.Body, &data); err != nil {
		return
	}

	if cur, _ := s.actions.GetProperty(ctx, id); cur != nil && cur.GetID() == id {
		resCode = http.StatusOK
	}

	property := data.ToProperty()
	property.ID = id

	if _, err := s.actions.StoreProperty(ctx, property); err != nil {
		errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	jsonResponse(w, resCode, oapi.NewStorePropertyRes(property),
		Header{"Location", "/property/" + property.ID})
}
func (s *Server) GetPropertyById(w http.ResponseWriter, r *http.Request, propertyID string) {
	ctx := r.Context()
	property, err := s.actions.GetProperty(ctx, propertyID)
	if err != nil {
		switch {
		case errors.Is(err, internal.ErrEntityNotFound):
			errorResponse(w, http.StatusNotFound, err.Error())
		case errors.Is(err, internal.ErrInternal):
			errorResponse(w, http.StatusInternalServerError, err.Error())
		default:
			errorResponse(w, http.StatusBadRequest, err.Error())
		}
		return
	}
	jsonResponse(w, http.StatusOK, oapi.NewGetPropertyRes(*property))
}
func (s *Server) ListProperties(w http.ResponseWriter, r *http.Request, params oapi.ListPropertiesParams) {
	var (
		ctx = r.Context()
		f   = params.ToFilter()
	)
	propList, err := s.actions.ListProperties(ctx, f)
	if err != nil {
		s.logError(err)
		errorResponse(w, http.StatusInternalServerError, "Error fetching list")
		return
	}
	jsonResponse(w, http.StatusOK, oapi.NewListPropertiesRes(propList...))
}
func (s *Server) DeleteProperty(w http.ResponseWriter, r *http.Request, propertyID string) {
	ctx := r.Context()
	if err := s.actions.RemoveProperty(ctx, propertyID); err != nil {
		switch {
		case errors.Is(err, internal.ErrInternal):
			errorResponse(w, http.StatusInternalServerError, err.Error())
			return
		case errors.Is(err, internal.ErrEntityNotFound):
		// what should a restful DELETE endpoint do
		// when the resource does not exist?
		// for now, I vote nothing, they client wants it gone,
		// and it isn't there ... so client should be happy
		default:
			errorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
	}
	w.WriteHeader(http.StatusNoContent)
}

func NewServer(acts actions.Actions) *Server {
	server := Server{
		Conf:    config.GetConfig(),
		actions: acts,
	}
	return &server
}
func (s *Server) WithConfig(conf config.Config) *Server {
	s2 := *s
	s2.Conf = conf
	return &s2
}
func (s *Server) Handler() http.Handler {
	router := chi.NewRouter()
	router.Group(func(r chi.Router) {
		r.Get("/health", s.okHandler)
		r.Get("/health/ready", s.okHandler)
		r.Get("/health/live", s.okHandler)
	})
	oapi.HandlerWithOptions(s, oapi.ChiServerOptions{
		BaseRouter: router,
		Middlewares: []oapi.MiddlewareFunc{
			s.AuthMW,
		},
	})
	return router
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
func (s *Server) okHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}
func (s *Server) logError(err error) {
	log.Entry().Error(err)
}

func errorResponse(w http.ResponseWriter, code int, msg string) {
	jsonResponse(w, code, oapi.ErrorResponse{
		Error: oapi.Error{
			Code:    0, // TODO
			Message: msg,
			Type:    "", // TODO
		},
	})
}
func jsonResponse(w http.ResponseWriter, resCode int, data interface{}, headers ...Header) {
	jData, err := json.Marshal(data)
	if err != nil {
		log.WithError(err).Error("json.Encode response failed")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	for _, h := range headers {
		w.Header().Set(h.k, h.v)
	}
	w.WriteHeader(resCode)
	if _, err := w.Write(jData); err != nil {
		log.WithError(err).Error("w.Write failed")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
func decodeRequestData(w http.ResponseWriter, body io.Reader, data interface{}) error {
	err := json.NewDecoder(body).Decode(&data)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, "Request body was not valid json")
		return err
	}
	return nil
}
