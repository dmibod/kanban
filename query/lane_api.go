package query

import (
	"context"
	"github.com/dmibod/kanban/shared/tools/mux"
	"github.com/go-chi/render"
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
	laneService services.LaneService
	cardService services.CardService
	logger.Logger
}

// CreateLaneAPI creates API
func CreateLaneAPI(lane services.LaneService, card services.CardService, l logger.Logger) *LaneAPI {
	return &LaneAPI{
		laneService: lane,
		cardService: card,
		Logger:      l,
	}
}

// Routes install handlers
func (a *LaneAPI) Routes(router chi.Router) {
	router.Get("/", a.All)
	router.Get("/{LANEID}", a.Get)
	router.Get("/{LANEID}/card", a.GetCards)
	router.Get("/{LANEID}/lane", a.GetLanes)
}

// All lanes
func (a *LaneAPI) All(w http.ResponseWriter, r *http.Request) {
	op := handlers.All(a, &LaneGetMapper{}, a.Logger)
	handlers.Handle(w, r, op)
}

// Get lane
func (a *LaneAPI) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "LANEID")
	op := handlers.Get(id, a, &LaneGetMapper{}, a.Logger)
	handlers.Handle(w, r, op)
}

// GetByID implements handlers.GetService
func (a *LaneAPI) GetByID(ctx context.Context, id kernel.Id) (interface{}, error) {
	return a.laneService.GetByID(ctx, id)
}

// GetCards by lane
func (a *LaneAPI) GetCards(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "LANEID")
	models, err := a.cardService.GetByLaneID(r.Context(), kernel.Id(id))
	if err != nil {
		a.Errorln(err)
		mux.RenderError(w, http.StatusInternalServerError)
		return
	}
	mapper := CardGetMapper{}
	render.JSON(w, r, mapper.ModelsToPayload(models))
}

// GetLanes by lane
func (a *LaneAPI) GetLanes(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "LANEID")
	models, err := a.laneService.GetByLaneID(r.Context(), kernel.Id(id))
	if err != nil {
		a.Errorln(err)
		mux.RenderError(w, http.StatusInternalServerError)
		return
	}
	mapper := LaneGetMapper{}
	render.JSON(w, r, mapper.ModelsToPayload(models))
}

// GetAll implements handlers.AllService
func (a *LaneAPI) GetAll(ctx context.Context) ([]interface{}, error) {
	models, err := a.laneService.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	mapper := LaneGetMapper{}
	return mapper.ModelsToPayload(models), nil
}

// LaneGetMapper mapper
type LaneGetMapper struct {
}

// ModelToPayload mapping
func (LaneGetMapper) ModelToPayload(m interface{}) interface{} {
	model := m.(*services.LaneModel)
	return &Lane{
		ID:     string(model.ID),
		Name:   model.Name,
		Type:   model.Type,
		Layout: model.Layout,
	}
}

// ModelsToPayload mapping
func (m LaneGetMapper) ModelsToPayload(models []*services.LaneModel) []interface{} {
	items := []interface{}{}
	for _, model := range models {
		items = append(items, m.ModelToPayload(model))
	}
	return items
}
