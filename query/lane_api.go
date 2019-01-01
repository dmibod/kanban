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
	router.Get("/{LANEID}", a.Get)
	router.Get("/{LANEID}/card", a.GetCards)
}

// Get lane
func (a *LaneAPI) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "LANEID")
	op := handlers.Get(id, a, &laneGetMapper{}, a.Logger)
	handlers.Handle(w, r, op)
}

// GetByID implements handlers.GetService
func (a *LaneAPI) GetByID(ctx context.Context, id kernel.Id) (interface{}, error) {
	return a.laneService.GetByID(ctx, id)
}

// GetCards by lane
func (a *LaneAPI) GetCards(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "LANEID")
	cards, err := a.cardService.GetByLaneID(r.Context(), kernel.Id(id))
	if err != nil {
		a.Errorln(err)
		mux.RenderError(w, http.StatusInternalServerError)
		return
	}
	render.JSON(w, r, cards)
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
