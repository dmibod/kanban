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

// LaneAPI dependencies
type LaneAPI struct {
	services.LaneService
	services.CardService
	logger.Logger
}

// CreateLaneAPI creates API
func CreateLaneAPI(lane services.LaneService, card services.CardService, l logger.Logger) *LaneAPI {
	return &LaneAPI{
		LaneService: lane,
		CardService: card,
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
func (a *LaneAPI) GetByID(ctx context.Context, id kernel.ID) (interface{}, error) {
	return a.LaneService.GetByID(ctx, id)
}

// GetCards by lane
func (a *LaneAPI) GetCards(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "LANEID")
	models, err := a.CardService.GetByLaneID(r.Context(), kernel.ID(id))
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
	models, err := a.LaneService.GetByLaneID(r.Context(), kernel.ID(id))
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
	models, err := a.LaneService.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	mapper := LaneGetMapper{}
	return mapper.ModelsToPayload(models), nil
}
