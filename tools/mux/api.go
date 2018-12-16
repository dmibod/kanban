package mux

import (
	"net/http"
)

type Mux interface {
	Get(string, http.Handler)
	Post(string, http.Handler)
	Handle(string, string, http.Handler)
}
