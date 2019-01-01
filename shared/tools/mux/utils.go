package mux

import (
	"encoding/json"
	"net/http"
)

// ParseJSON request
func ParseJSON(r *http.Request, payload interface{}) error {
	return json.NewDecoder(r.Body).Decode(payload)
}

// RenderJSON response
func RenderJSON(w http.ResponseWriter, payload interface{}) {
	json.NewEncoder(w).Encode(payload)
}

// RenderError response
func RenderError(w http.ResponseWriter, code int) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	resp := struct {
		Message string `json:"message"`
		Success bool   `json:"success"`
	}{http.StatusText(code), false}
	RenderJSON(w, resp)
}
