package common

import (
	"encoding/json"
	"net/http"

	"github.com/gold-kou/go-housework/app/model/schemamodel"
	schema "github.com/gorilla/Schema"
	log "github.com/sirupsen/logrus"
)

var decoder = schema.NewDecoder()

// ResponseBadRequest sets response as bad request
func ResponseBadRequest(w http.ResponseWriter, message string, errors ...map[string]string) {
	res := schemamodel.ResponseBadRequest{
		Message: message,
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusBadRequest)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Error(err)
		panic(err.Error())
	}
}

// ResponseUnauthorized sets response as unauthorized
func ResponseUnauthorized(w http.ResponseWriter, message string) {
	res := schemamodel.ResponseUnauthorized{
		Message: message,
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusUnauthorized)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Error(err)
		panic(err.Error())
	}
}

// ResponseInternalServerError sets response as internal server error
func ResponseInternalServerError(w http.ResponseWriter, message string) {
	res := schemamodel.ResponseInternalServerError{
		Message: message,
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusInternalServerError)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Error(err)
		panic(err.Error())
	}
}
