package mux

import (
	"net/http"
)

type Mux interface {
	Any(string, http.Handler)
	Get(string, http.Handler)
	Post(string, http.Handler)
	Handle(string, string, http.Handler)
}
