package update

import (
	"context"
	"net/http"

	"github.com/dmibod/kanban/shared/handlers"

	"github.com/dmibod/kanban/shared/kernel"
	"github.com/dmibod/kanban/shared/services"
	"github.com/go-chi/chi"

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
	router.Post("/", a.CreateLane)
	router.Put("/{LANEID}", a.UpdateLane)
	router.Delete("/{LANEID}", a.RemoveLane)
}

// CreateLane handler
func (a *LaneAPI) CreateLane(w http.ResponseWriter, r *http.Request) {
	op := handlers.Create(&Lane{}, a, &laneCreateMapper{}, a.Logger)
	handlers.Handle(w, r, op)
}

// UpdateLane handler
func (a *LaneAPI) UpdateLane(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "LANEID")
	op := handlers.Update(&Lane{ID: id}, a, &laneUpdateMapper{}, a.Logger)
	handlers.Handle(w, r, op)
}

// RemoveLane handler
func (a *LaneAPI) RemoveLane(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "LANEID")
	op := handlers.Remove(id, a.LaneService, a.Logger)
	handlers.Handle(w, r, op)
}

// Create implements handlers.CreateService
func (a *LaneAPI) Create(ctx context.Context, model interface{}) (interface{}, error) {
	return a.LaneService.Create(ctx, model.(*services.LanePayload))
}

// Update implements handlers.UpdateService
func (a *LaneAPI) Update(ctx context.Context, model interface{}) (interface{}, error) {
	return a.LaneService.Update(ctx, model.(*services.LaneModel))
}

type laneCreateMapper struct {
}

// PayloadToModel mapping
func (laneCreateMapper) PayloadToModel(p interface{}) interface{} {
	payload := p.(*Lane)
	return &services.LanePayload{
		Name:   payload.Name,
		Type:   payload.Type,
		Layout: payload.Layout,
	}
}

// ModelToPayload mapping
func (laneCreateMapper) ModelToPayload(m interface{}) interface{} {
	model := m.(*services.LaneModel)
	return &Lane{
		ID:     string(model.ID),
		Name:   model.Name,
		Type:   model.Type,
		Layout: model.Layout,
	}
}

type laneUpdateMapper struct {
}

// PayloadToModel mapping
func (laneUpdateMapper) PayloadToModel(p interface{}) interface{} {
	payload := p.(*Lane)
	return &services.LaneModel{
		ID:     kernel.ID(payload.ID),
		Name:   payload.Name,
		Type:   payload.Type,
		Layout: payload.Layout,
	}
}

// ModelToPayload mapping
func (laneUpdateMapper) ModelToPayload(m interface{}) interface{} {
	model := m.(*services.LaneModel)
	return &Lane{
		ID:     string(model.ID),
		Name:   model.Name,
		Type:   model.Type,
		Layout: model.Layout,
	}
}
