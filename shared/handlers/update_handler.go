package handlers

import (
	"context"
	"net/http"

	"github.com/dmibod/kanban/shared/tools/logger"
	"github.com/dmibod/kanban/shared/tools/mux"
	"github.com/go-chi/render"
)

// Update operation
func Update(payload interface{}, service UpdateService, mapper MapService, l logger.Logger) Operation {
	return &updateOperation{
		payload: payload,
		service: service,
		mapper:  mapper,
		Logger:  l,
	}
}

// UpdateService interface
type UpdateService interface {
	Update(context.Context, interface{}) (interface{}, error)
}

type updateOperation struct {
	logger.Logger
	payload interface{}
	service UpdateService
	mapper  MapService
}

// Execute update
func (o *updateOperation) Execute(w http.ResponseWriter, r *http.Request) {
	err := mux.ParseJSON(r, o.payload)
	if err != nil {
		o.Errorln(err)
		mux.RenderError(w, http.StatusBadRequest)
		return
	}

	model, err := o.service.Update(r.Context(), o.mapper.PayloadToModel(o.payload))
	if err != nil {
		o.Errorln(err)
		mux.RenderError(w, http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, o.mapper.ModelToPayload(model))
}
