package lane

import (
	"github.com/dmibod/kanban/shared/kernel"
	"net/http"

	"github.com/dmibod/kanban/shared/handlers"

	"github.com/dmibod/kanban/shared/services/lane"
	"github.com/go-chi/chi"

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
	router.Post("/{BOARDID}/lanes", a.CreateLane)
}

// CreateLane handler
func (a *API) CreateLane(w http.ResponseWriter, r *http.Request) {
	op := handlers.Create(&Lane{}, a, &laneCreateMapper{}, a.Logger)
	handlers.Handle(w, r, op)
}

// Create implements handlers.CreateService
func (a *API) Create(r *http.Request, model interface{}) (interface{}, error) {
	boardID := chi.URLParam(r, "BOARDID")
	laneID, err := a.Service.Create(r.Context(), kernel.ID(boardID), model.(*lane.CreateModel))
	if err != nil {
		return nil, err
	}
	return a.Service.GetByID(r.Context(), laneID.WithSet(kernel.ID(boardID)))
}
