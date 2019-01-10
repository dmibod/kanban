package handlers

import (
	"context"
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
	Create(context.Context, interface{}) (interface{}, error)
}

type createOperation struct {
	logger.Logger
	payload interface{}
	service CreateService
	mapper  MapService
}

// Execute create
func (o *createOperation) Execute(w http.ResponseWriter, r *http.Request) {
	err := mux.ParseJSON(r, o.payload)
	if err != nil {
		o.Errorln(err)
		mux.RenderError(w, http.StatusBadRequest)
		return
	}

	model, err := o.service.Create(r.Context(), o.mapper.PayloadToModel(o.payload))
	if err != nil {
		o.Errorln(err)
		mux.RenderError(w, http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, o.mapper.ModelToPayload(model))
}
