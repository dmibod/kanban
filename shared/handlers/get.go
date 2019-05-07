package handlers

import (
	"net/http"

	"github.com/dmibod/kanban/shared/tools/logger"
	"github.com/dmibod/kanban/shared/tools/mux"
	"github.com/go-chi/render"
)

// Get operation
func Get(service GetService, mapper ModelMapper, l logger.Logger) Operation {
	return &getOperation{
		service: service,
		mapper:  mapper,
		Logger:  l,
	}
}

// GetService interface
type GetService interface {
	GetOne(*http.Request) (interface{}, error)
}

type getOperation struct {
	logger.Logger
	service GetService
	mapper  ModelMapper
}

// Execute get
func (o *getOperation) Execute(w http.ResponseWriter, r *http.Request) {
	if model, err := o.service.GetOne(r); err != nil {
		o.Errorln(err)
		mux.RenderError(w, http.StatusNotFound)
	} else {
		render.JSON(w, r, o.mapper.ModelToPayload(model))
	}
}
