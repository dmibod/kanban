package lane

import (
	"context"
	"net/http"

	cardapi "github.com/dmibod/kanban/query/card"

	"github.com/dmibod/kanban/shared/services/card"
	"github.com/dmibod/kanban/shared/services/lane"

	"github.com/dmibod/kanban/shared/tools/mux"
	"github.com/go-chi/render"

	"github.com/dmibod/kanban/shared/handlers"

	"github.com/go-chi/chi"

	"github.com/dmibod/kanban/shared/kernel"
	"github.com/dmibod/kanban/shared/tools/logger"
)

// API dependencies
type API struct {
	laneService lane.Service
	cardService card.Service
	logger.Logger
}

// CreateAPI creates API
func CreateAPI(ls lane.Service, cs card.Service, l logger.Logger) *API {
	return &API{
		laneService: ls,
		cardService: cs,
		Logger:      l,
	}
}

// Routes install handlers
func (a *API) Routes(router chi.Router) {
	router.Get("/", a.All)
	router.Get("/{LANEID}", a.Get)
	router.Get("/{LANEID}/card", a.GetCards)
	router.Get("/{LANEID}/lane", a.GetLanes)
}

// All lanes
func (a *API) All(w http.ResponseWriter, r *http.Request) {
	op := handlers.All(a, &ListModelMapper{}, a.Logger)
	handlers.Handle(w, r, op)
}

// Get lane
func (a *API) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "LANEID")
	op := handlers.Get(id, a, &ModelMapper{}, a.Logger)
	handlers.Handle(w, r, op)
}

// GetByID implements handlers.GetService
func (a *API) GetByID(ctx context.Context, id kernel.ID) (interface{}, error) {
	return a.laneService.GetByID(ctx, id)
}

// GetCards by lane
func (a *API) GetCards(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "LANEID")
	if models, err := a.cardService.GetByLaneID(r.Context(), kernel.ID(id)); err == nil {
		mapper := cardapi.CardGetMapper{}
		render.JSON(w, r, mapper.ModelsToPayload(models))
	} else {
		a.Errorln(err)
		mux.RenderError(w, http.StatusInternalServerError)
	}
}

// GetLanes by lane
func (a *API) GetLanes(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "LANEID")
	if models, err := a.laneService.GetByLaneID(r.Context(), kernel.ID(id)); err == nil {
		mapper := ListModelMapper{}
		render.JSON(w, r, mapper.ModelsToPayload(models))
	} else {
		a.Errorln(err)
		mux.RenderError(w, http.StatusInternalServerError)
	}
}

// GetAll implements handlers.AllService
func (a *API) GetAll(ctx context.Context) ([]interface{}, error) {
	models, err := a.laneService.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	mapper := ListModelMapper{}
	return mapper.ModelsToPayload(models), nil
}
