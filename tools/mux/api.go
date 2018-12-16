package mux

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// Mux abstracts mux
type Mux interface {
	All(string, http.Handler)
	Get(string, http.Handler)
	Post(string, http.Handler)
}

// ApiFunc func to server api request
type ApiFunc func(interface{}) (interface{}, error)

// FactoryFunc func to instantiate api request 
type FactoryFunc func() interface{}

// ApiHandler type to serve as handler
type ApiHandler struct {
	h ApiFunc
	f FactoryFunc
}

func (h *ApiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	req := h.f()
	reqErr := JsonRequest(r, req)
	if reqErr != nil {
		ErrorResponse(w, http.StatusInternalServerError)
	}
	res, resErr := h.h(req)
	if resErr != nil {
		ErrorResponse(w, http.StatusInternalServerError)
	}
	JsonResponse(w, res)
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
