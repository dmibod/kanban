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

// ApiHandler type to serve as handler
type ApiHandler interface {
	ParseRequest(*http.Request) (interface{}, error)
	Handle(interface{}) (interface{}, error)
}

// ApiHandleFunc returns HandleFunc from ApiHandler
func ApiHandleFunc(h ApiHandler) http.HandlerFunc {
   return func(w http.ResponseWriter, r *http.Request){
		req, reqErr := h.ParseRequest(r)
		if reqErr != nil {
			ErrorResponse(w, http.StatusInternalServerError)
		}
		res, resErr := h.Handle(req)
		if resErr != nil {
			ErrorResponse(w, http.StatusInternalServerError)
		}
		JsonResponse(w, res)
	 }
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
