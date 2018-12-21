package mux

import (
	"encoding/json"
	"net/http"
)

// JsonRequest - parses request as json
func JsonRequest(r *http.Request, payload interface{}) error {
	return json.NewDecoder(r.Body).Decode(payload)
}

// JsonResponse - builds json response
func JsonResponse(w http.ResponseWriter, payload interface{}) {
	json.NewEncoder(w).Encode(payload)
}

// ErrorResponse - builds error response
func ErrorResponse(w http.ResponseWriter, code int) {
	http.Error(w, http.StatusText(code), code)
}
