package update

import (
	"context"
	"net/http"

	"github.com/dmibod/kanban/shared/handlers"

	"github.com/dmibod/kanban/shared/services"
	"github.com/go-chi/chi"

	"github.com/dmibod/kanban/shared/tools/logger"
)

// Lane model
type Lane struct {
	ID          string `json:"id,omitempty"`
	Type        string `json:"type"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Layout      string `json:"layout"`
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
	router.Post("/", a.CreateLane)
}

// CreateLane handler
func (a *LaneAPI) CreateLane(w http.ResponseWriter, r *http.Request) {
	op := handlers.Create(&Lane{}, a, &laneCreateMapper{}, a.Logger)
	handlers.Handle(w, r, op)
}

// Create implements handlers.CreateService
func (a *LaneAPI) Create(ctx context.Context, model interface{}) (interface{}, error) {
	return a.LaneService.Create(ctx, model.(*services.LanePayload))
}

type laneCreateMapper struct {
}

// PayloadToModel mapping
func (laneCreateMapper) PayloadToModel(p interface{}) interface{} {
	payload := p.(*Lane)
	return &services.LanePayload{
		Type:        payload.Type,
		Name:        payload.Name,
		Description: payload.Description,
		Layout:      payload.Layout,
	}
}

// ModelToPayload mapping
func (laneCreateMapper) ModelToPayload(m interface{}) interface{} {
	model := m.(*services.LaneModel)
	return &Lane{
		ID:          string(model.ID),
		Type:        model.Type,
		Name:        model.Name,
		Description: model.Description,
		Layout:      model.Layout,
	}
}
