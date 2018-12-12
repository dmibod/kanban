package mux

import (
	"net/http"
)
type Mux interface {
	Handle(string, http.Handler)
}