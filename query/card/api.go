package card

import (
	"context"
	"net/http"

	"github.com/dmibod/kanban/shared/handlers"

	"github.com/dmibod/kanban/shared/services/card"
	"github.com/go-chi/chi"

	"github.com/dmibod/kanban/shared/kernel"
	"github.com/dmibod/kanban/shared/tools/logger"
)

// API dependencies
type API struct {
	card.Service
	logger.Logger
}

// CreateAPI creates new instance of API
func CreateAPI(s card.Service, l logger.Logger) *API {
	return &API{
		Service: s,
		Logger:  l,
	}
}

// Routes install API handlers
func (a *API) Routes(router chi.Router) {
	router.Get("/", a.All)
	router.Get("/{CARDID}", a.Get)
}

// All cards
func (a *API) All(w http.ResponseWriter, r *http.Request) {
	op := handlers.All(a, &CardGetMapper{}, a.Logger)
	handlers.Handle(w, r, op)
}

// Get card by id
func (a *API) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "CARDID")
	op := handlers.Get(id, a, &CardGetMapper{}, a.Logger)
	handlers.Handle(w, r, op)
}

// GetAll implements handlers.AllService
func (a *API) GetAll(ctx context.Context) ([]interface{}, error) {
	if models, err := a.Service.GetAll(ctx); err != nil {
		return nil, err
	} else {
		mapper := CardGetMapper{}
		return mapper.ModelsToPayload(models), nil
	}
}

// GetByID implements handlers.GetService
func (a *API) GetByID(ctx context.Context, id kernel.ID) (interface{}, error) {
	return a.Service.GetByID(ctx, id)
}
