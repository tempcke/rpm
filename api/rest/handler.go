package rest

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/tempcke/rpm/pkg/log"
)

func decodeRequestData(w http.ResponseWriter, body io.Reader, data interface{}) error {
	err := json.NewDecoder(body).Decode(&data)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, "Request body was not valid json")
		return err
	}
	return nil
}

func jsonResponse(w http.ResponseWriter, resCode int, data interface{}) {
	jData, err := json.Marshal(data)
	if err != nil {
		log.WithError(err).Error("json.Encode response failed")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resCode)
	if _, err := w.Write(jData); err != nil {
		log.WithError(err).Error("w.Write failed")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func errorResponse(w http.ResponseWriter, code int, msg string) {
	jsonResponse(w, code, ErrorResponse{Error: msg})
}
