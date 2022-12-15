package rest

import (
	"net/http"

	"github.com/go-chi/chi"
	log "github.com/sirupsen/logrus"
	"github.com/tempcke/rpm/usecase"
)

func addProperty(propRepo usecase.PropertyRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		data := PropertyModel{}
		if err := decodeRequestData(w, r.Body, &data); err != nil {
			return
		}

		uc := usecase.NewAddProperty(propRepo)
		property := propRepo.NewProperty(
			data.Street, data.City, data.State, data.Zip,
		)

		if err := uc.Execute(ctx, property); err != nil {
			errorResponse(w, http.StatusBadRequest, "Missing or invalid fields")
			return
		}

		w.WriteHeader(http.StatusCreated)
		jsonResponse(w, NewPropertyModel(property))
	}
}

func listProperties(propRepo usecase.PropertyReader) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		uc := usecase.NewListProperties(propRepo)
		propList, err := uc.Execute(ctx)
		if err != nil {
			log.Error(err)
			errorResponse(w, http.StatusNotFound, "Error fetching list")
			return
		}
		jsonResponse(w, NewPropertyListModel(propList...))
	}
}

func getProperty(propRepo usecase.PropertyReader) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		propertyID := chi.URLParam(r, "propertyID")
		uc := usecase.NewGetProperty(propRepo)
		property, err := uc.Execute(ctx, propertyID)
		if err != nil {
			errorResponse(w, http.StatusNotFound, "propertyId not found")
			return
		}
		jsonResponse(w, NewPropertyModel(property))
	}
}

func deleteProperty(propRepo usecase.PropertyRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		propertyID := chi.URLParam(r, "propertyID")
		uc := usecase.NewDeleteProperty(propRepo)
		err := uc.Execute(ctx, propertyID)
		if err != nil {
			// what should a restful DELETE endpoint do
			// when the resource does not exist?
			// for now, I vote nothing, they client wants it gone,
			// and it isn't there ... so client should be happy

			w.WriteHeader(http.StatusNoContent)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
