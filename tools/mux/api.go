package mux

import (
	"encoding/json"
	"net/http"
)

// Mux abstracts mux
type Mux interface {
	All(string, http.Handler)
	Get(string, http.Handler)
	Post(string, http.Handler)
}

// Json - builds json response
func Json(w http.ResponseWriter, payload interface{}) {
	enc := json.NewEncoder(w)
	enc.Encode(payload)
}

// Error - builds error response
func Error(w http.ResponseWriter, code int) {
	http.Error(w, http.StatusText(code), code)
}
