package query

import (
	"context"
	"net/http"

	"github.com/dmibod/kanban/shared/handlers"

	"github.com/dmibod/kanban/shared/services"
	"github.com/go-chi/chi"

	"github.com/dmibod/kanban/shared/kernel"
	"github.com/dmibod/kanban/shared/tools/logger"
)

// Lane model
type Lane struct {
	ID     string `json:"id,omitempty"`
	Name   string `json:"name,omitempty"`
	Type   string `json:"type"`
	Layout string `json:"layout"`
}

// LaneAPI dependencies
type LaneAPI struct {
	services.LaneService
	logger.Logger
}

// CreateLaneAPI creates API
func CreateLaneAPI(s services.LaneService, l logger.Logger) *LaneAPI {
	return &LaneAPI{
		LaneService: s,
		Logger:      l,
	}
}

// Routes install handlers
func (a *LaneAPI) Routes(router chi.Router) {
	router.Get("/{LANEID}", a.Get)
}

// Get lane
func (a *LaneAPI) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "LANEID")
	op := handlers.Get(id, a, &laneGetMapper{}, a.Logger)
	handlers.Handle(w, r, op)
}

// GetByID implements handlers.GetService
func (a *LaneAPI) GetByID(ctx context.Context, id kernel.Id) (interface{}, error) {
	return a.LaneService.GetByID(ctx, id)
}

type laneGetMapper struct {
}

// ModelToPayload mapping
func (laneGetMapper) ModelToPayload(m interface{}) interface{} {
	model := m.(*services.LaneModel)
	return &Lane{
		ID:     string(model.ID),
		Name:   model.Name,
		Type:   model.Type,
		Layout: model.Layout,
	}
}
