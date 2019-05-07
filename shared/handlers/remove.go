package handlers

import (
	"context"
	"net/http"

	"github.com/dmibod/kanban/shared/kernel"
	"github.com/dmibod/kanban/shared/tools/logger"
	"github.com/dmibod/kanban/shared/tools/mux"
	"github.com/go-chi/render"
)

// Remove operation
func Remove(id string, service RemoveService, l logger.Logger) Operation {
	return &removeOperation{
		id:      id,
		service: service,
		Logger:  l,
	}
}

// RemoveService interface
type RemoveService interface {
	Remove(context.Context, kernel.ID) error
}

type removeOperation struct {
	logger.Logger
	id      string
	service RemoveService
}

// Execute remove
func (o *removeOperation) Execute(w http.ResponseWriter, r *http.Request) {
	if err := o.service.Remove(r.Context(), kernel.ID(o.id)); err != nil {
		o.Errorln(err)
		mux.RenderError(w, http.StatusInternalServerError)
		return
	}

	resp := struct {
		ID      string `json:"id"`
		Success bool   `json:"success"`
	}{o.id, true}

	render.JSON(w, r, resp)
}
