package handlers

import (
	"net/http"

	"github.com/dmibod/kanban/shared/tools/logger"
	"github.com/dmibod/kanban/shared/tools/mux"
	"github.com/go-chi/render"
)

// Create operation
func Create(payload interface{}, service CreateService, mapper MapService, l logger.Logger) Operation {
	return &createOperation{
		payload: payload,
		service: service,
		mapper:  mapper,
		Logger:  l,
	}
}

// CreateService interface
type CreateService interface {
	Create(*http.Request, interface{}) (interface{}, error)
}

type createOperation struct {
	logger.Logger
	payload interface{}
	service CreateService
	mapper  MapService
}

// Execute create
func (o *createOperation) Execute(w http.ResponseWriter, r *http.Request) {
	if err := mux.ParseJSON(r, o.payload); err != nil {
		o.Errorln(err)
		mux.RenderError(w, http.StatusBadRequest)
		return
	}

	if model, err := o.service.Create(r, o.mapper.PayloadToModel(o.payload)); err != nil {
		o.Errorln(err)
		mux.RenderError(w, http.StatusInternalServerError)
	} else {
		render.JSON(w, r, o.mapper.ModelToPayload(model))
	}
}
