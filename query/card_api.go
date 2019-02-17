package query

import (
	"context"
	"net/http"

	"github.com/dmibod/kanban/shared/handlers"

	"github.com/dmibod/kanban/shared/services"
	"github.com/go-chi/chi"

	"github.com/dmibod/kanban/shared/kernel"
	"github.com/dmibod/kanban/shared/tools/logger"
)

// CardAPI dependencies
type CardAPI struct {
	services.CardService
	logger.Logger
}

// CreateCardAPI creates new instance of API
func CreateCardAPI(s services.CardService, l logger.Logger) *CardAPI {
	return &CardAPI{
		CardService: s,
		Logger:      l,
	}
}

// Routes install API handlers
func (a *CardAPI) Routes(router chi.Router) {
	router.Get("/", a.All)
	router.Get("/{CARDID}", a.Get)
}

// All cards
func (a *CardAPI) All(w http.ResponseWriter, r *http.Request) {
	op := handlers.All(a, &CardGetMapper{}, a.Logger)
	handlers.Handle(w, r, op)
}

// Get card by id
func (a *CardAPI) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "CARDID")
	op := handlers.Get(id, a, &CardGetMapper{}, a.Logger)
	handlers.Handle(w, r, op)
}

// GetAll implements handlers.AllService
func (a *CardAPI) GetAll(ctx context.Context) ([]interface{}, error) {
	models, err := a.CardService.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	mapper := CardGetMapper{}
	return mapper.ModelsToPayload(models), nil
}

// GetByID implements handlers.GetService
func (a *CardAPI) GetByID(ctx context.Context, id kernel.ID) (interface{}, error) {
	return a.CardService.GetByID(ctx, id)
}