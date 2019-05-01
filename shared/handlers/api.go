package handlers

import (
	"net/http"
)

// Operation interface
type Operation interface {
	Execute(w http.ResponseWriter, r *http.Request)
}

// Handler definition
type Handler struct {
	Operation
}

// Handle operation
func (h Handler) Handle(w http.ResponseWriter, r *http.Request) {
	h.Operation.Execute(w, r)
}

// Handle operation
func Handle(w http.ResponseWriter, r *http.Request, o Operation) {
	h := Handler{Operation: o}
	h.Handle(w, r)
}

// PayloadMapper interface
type PayloadMapper interface {
	PayloadToModel(interface{}) interface{}
}

// ModelMapper interface
type ModelMapper interface {
	ModelToPayload(interface{}) interface{}
}

// MapService interface
type MapService interface {
	PayloadMapper
	ModelMapper
}
