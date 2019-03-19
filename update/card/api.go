package card

import (
	"context"
	"net/http"

	"github.com/dmibod/kanban/shared/handlers"

	"github.com/dmibod/kanban/shared/services/card"
	"github.com/go-chi/chi"

	"github.com/dmibod/kanban/shared/tools/logger"
)

// API dependencies
type API struct {
	service card.Service
	logger.Logger
}

// CreateAPI creates API
func CreateAPI(s card.Service, l logger.Logger) *API {
	return &API{
		service: s,
		Logger:  l,
	}
}

// Routes install handlers
func (a *API) Routes(router chi.Router) {
	router.Post("/", a.CreateCard)
}

// CreateCard handler
func (a *API) CreateCard(w http.ResponseWriter, r *http.Request) {
	op := handlers.Create(&Card{}, a, &cardCreateMapper{}, a.Logger)
	handlers.Handle(w, r, op)
}

// Create implements handlers.CreateService
func (a *API) Create(ctx context.Context, model interface{}) (interface{}, error) {
	id, err := a.service.Create(ctx, model.(*card.CreateModel))
	if err != nil {
		return nil, err
	}
	return a.service.GetByID(ctx, id)
}
