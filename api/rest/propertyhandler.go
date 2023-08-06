package rest

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
	"github.com/tempcke/rpm/actions"
	"github.com/tempcke/rpm/entity"
	"github.com/tempcke/rpm/internal"
)

type propertyHandler struct {
	actions actions.Actions
}

func (h propertyHandler) addProperty(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	data := PropertyModel{}
	if err := decodeRequestData(w, r.Body, &data); err != nil {
		return
	}

	property := entity.NewProperty(
		data.Street, data.City, data.State, data.Zip,
	)

	if _, err := h.actions.StoreProperty(ctx, property); err != nil {
		errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Location", "/property/"+property.ID)
	jsonResponse(w, http.StatusCreated, NewPropertyModel(property))
}
func (h propertyHandler) storeProperty(w http.ResponseWriter, r *http.Request) {
	var (
		ctx        = r.Context()
		propertyID = chi.URLParam(r, "propertyID")
		resCode    = http.StatusCreated
	)

	if curProp, _ := h.actions.GetProperty(ctx, propertyID); curProp != nil && curProp.ID == propertyID {
		resCode = http.StatusOK
	}

	data := PropertyModel{}
	if err := decodeRequestData(w, r.Body, &data); err != nil {
		return
	}

	property := entity.NewProperty(data.Street, data.City, data.State, data.Zip)
	property.ID = propertyID

	if _, err := h.actions.StoreProperty(ctx, property); err != nil {
		errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// we want to reply with 201 when it is new and 200 when updated
	jsonResponse(w, resCode, NewPropertyModel(property))
	w.Header().Set("Location", "/property/"+property.ID)
}
func (h propertyHandler) listProperties(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	propList, err := h.actions.ListProperties(ctx)
	if err != nil {
		log.Error(err)
		errorResponse(w, http.StatusNotFound, "Error fetching list")
		return
	}
	jsonResponse(w, http.StatusOK, NewPropertyListModel(propList...))
}
func (h propertyHandler) getProperty(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	propertyID := chi.URLParam(r, "propertyID")
	property, err := h.actions.GetProperty(ctx, propertyID)
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
	jsonResponse(w, http.StatusOK, NewPropertyModel(*property))
}
func (h propertyHandler) deleteProperty(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	propertyID := chi.URLParam(r, "propertyID")
	if err := h.actions.RemoveProperty(ctx, propertyID); err != nil {
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
