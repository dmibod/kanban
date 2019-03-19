package lane

import (
	"context"
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
	router.Post("/", a.CreateLane)
}

// CreateLane handler
func (a *API) CreateLane(w http.ResponseWriter, r *http.Request) {
	op := handlers.Create(&Lane{}, a, &laneCreateMapper{}, a.Logger)
	handlers.Handle(w, r, op)
}

// Create implements handlers.CreateService
func (a *API) Create(ctx context.Context, model interface{}) (interface{}, error) {
	return a.Service.Create(ctx, model.(*lane.CreateModel))
}
