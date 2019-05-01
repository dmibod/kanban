package handlers

import (
	"net/http"

	"github.com/dmibod/kanban/shared/tools/logger"
	"github.com/dmibod/kanban/shared/tools/mux"
	"github.com/go-chi/render"
)

// List operation
func List(service ListService, mapper ModelMapper, l logger.Logger) Operation {
	return &listOperation{
		service: service,
		mapper:  mapper,
		Logger:  l,
	}
}

// ListService interface
type ListService interface {
	GetList(*http.Request) ([]interface{}, error)
}

type listOperation struct {
	logger.Logger
	service ListService
	mapper  ModelMapper
}

// Execute get
func (o *listOperation) Execute(w http.ResponseWriter, r *http.Request) {
	if models, err := o.service.GetList(r); err != nil {
		o.Errorln(err)
		mux.RenderError(w, http.StatusInternalServerError)
	} else {
		items := make([]interface{}, len(models))
		for i, model := range models {
			items[i] = o.mapper.ModelToPayload(model)
		}
		render.JSON(w, r, items)
	}
}
