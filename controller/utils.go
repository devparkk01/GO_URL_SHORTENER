package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

const (
	errorPathParamFailedParseNotFound = "A required path parameter was not found: '%s'"
	contentType                       = "Content-Type"
	applicationJson                   = "application/json"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func ServerResponse(w http.ResponseWriter, statusCode int, body interface{}) {
	SetHeader(w, contentType, applicationJson)
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(body)
}

func SetHeader(w http.ResponseWriter, key, value string) {
	w.Header().Set(key, value)
}

func ParsePathParam(r *http.Request, pathParam string) (string, error) {
	vars := mux.Vars(r)
	value, ok := vars[pathParam]
	if !ok {
		return "", errors.New(fmt.Sprintf(errorPathParamFailedParseNotFound, pathParam))
	}
	return value, nil
}
