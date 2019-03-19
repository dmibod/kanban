package handlers

import (
	"context"
	"net/http"

	"github.com/dmibod/kanban/shared/tools/logger"
	"github.com/dmibod/kanban/shared/tools/mux"
	"github.com/go-chi/render"
)

// All operation
func All(service AllService, mapper ModelMapper, l logger.Logger) Operation {
	return &allOperation{
		service: service,
		mapper:  mapper,
		Logger:  l,
	}
}

// AllService interface
type AllService interface {
	GetAll(context.Context) ([]interface{}, error)
}

type allOperation struct {
	logger.Logger
	service AllService
	mapper  ModelMapper
}

// Execute get
func (o *allOperation) Execute(w http.ResponseWriter, r *http.Request) {
	if models, err := o.service.GetAll(r.Context()); err != nil {
		o.Errorln(err)
		mux.RenderError(w, http.StatusInternalServerError)
	} else {
		items := []interface{}{}
		for _, model := range models {
			items = append(items, o.mapper.ModelToPayload(model))
		}
		render.JSON(w, r, items)
	}
}
