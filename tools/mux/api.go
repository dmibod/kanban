package mux

import (
	"io/ioutil"
	"encoding/json"
	"net/http"
)

// Mux abstracts mux
type Mux interface {
	All(string, http.Handler)
	Get(string, http.Handler)
	Post(string, http.Handler)
}

// JsonRequest - parses request as json
func JsonRequest(r *http.Request, payload interface{}) error {
	body, readErr := ioutil.ReadAll(r.Body)
	if readErr != nil {
		return readErr
	}

	jsonErr := json.Unmarshal(body, payload)
	if jsonErr != nil {
		return jsonErr
	}

	return nil
}

// JsonResponse - builds json response
func JsonResponse(w http.ResponseWriter, payload interface{}) {
	enc := json.NewEncoder(w)
	enc.Encode(payload)
}

// ErrorResponse - builds error response
func ErrorResponse(w http.ResponseWriter, code int) {
	http.Error(w, http.StatusText(code), code)
}
