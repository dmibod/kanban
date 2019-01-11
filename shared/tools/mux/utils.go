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
	http.Error(w, http.StatusText(code), code)
	/*
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		fmt.Fprintf(w, `{"success":"false","message":%q}`, http.StatusText(code))
	*/
}
