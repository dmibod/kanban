package handlers

import (
	"context"
	"github.com/dmibod/kanban/shared/kernel"
	"net/http"

	"github.com/dmibod/kanban/shared/tools/logger"
	"github.com/dmibod/kanban/shared/tools/mux"
	"github.com/go-chi/render"
)

// Get operation
func Get(id string, service GetService, mapper ModelMapper, l logger.Logger) Operation {
	return &getOperation{
		id:      id,
		service: service,
		mapper:  mapper,
		Logger:  l,
	}
}

// GetService interface
type GetService interface {
	GetByID(context.Context, kernel.ID) (interface{}, error)
}

type getOperation struct {
	logger.Logger
	id      string
	service GetService
	mapper  ModelMapper
}

// Execute get
func (o *getOperation) Execute(w http.ResponseWriter, r *http.Request) {
	if model, err := o.service.GetByID(r.Context(), kernel.ID(o.id)); err != nil {
		o.Errorln(err)
		mux.RenderError(w, http.StatusNotFound)
	} else {
		render.JSON(w, r, o.mapper.ModelToPayload(model))
	}
}
