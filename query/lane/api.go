package lane

import (
	"net/http"

	"github.com/dmibod/kanban/shared/services/lane"

	"github.com/dmibod/kanban/shared/handlers"

	"github.com/go-chi/chi"

	"github.com/dmibod/kanban/shared/kernel"
	"github.com/dmibod/kanban/shared/tools/logger"
)

// API dependencies
type API struct {
	lane.Service
	logger.Logger
}

// CreateAPI creates API
func CreateAPI(s lane.Service, l logger.Logger) *API {
	return &API{
		Service: s,
		Logger:  l,
	}
}

// Routes install handlers
func (a *API) Routes(router chi.Router) {
	router.Get("/{BOARDID}/lanes", a.List)
	router.Get("/{BOARDID}/lanes/{LANEID}", a.Get)
	router.Get("/{BOARDID}/lanes/{LANEID}/lanes", a.List)
}

// List lanes
func (a *API) List(w http.ResponseWriter, r *http.Request) {
	op := handlers.List(a, &ListModelMapper{}, a.Logger)
	handlers.Handle(w, r, op)
}

// GetList implements handlers.AllService
func (a *API) GetList(r *http.Request) ([]interface{}, error) {
	boardID := chi.URLParam(r, "BOARDID")
	laneID := chi.URLParam(r, "LANEID")
	var models []*lane.ListModel
	var err error
	if kernel.ID(laneID).IsValid() {
		if models, err = a.Service.GetByLaneID(r.Context(), kernel.ID(boardID).WithID(kernel.ID(laneID))); err != nil {
			return nil, err
		}
	} else {
		if models, err = a.Service.GetByBoardID(r.Context(), kernel.ID(boardID)); err != nil {
			return nil, err
		}
	}
	mapper := ListModelMapper{}
	return mapper.List(models), nil
}

// Get lane
func (a *API) Get(w http.ResponseWriter, r *http.Request) {
	op := handlers.Get(a, &ModelMapper{}, a.Logger)
	handlers.Handle(w, r, op)
}

// GetOne implements handlers.GetService
func (a *API) GetOne(r *http.Request) (interface{}, error) {
	boardID := chi.URLParam(r, "BOARDID")
	laneID := chi.URLParam(r, "LANEID")
	return a.Service.GetByID(r.Context(), kernel.ID(boardID).WithID(kernel.ID(laneID)))
}
