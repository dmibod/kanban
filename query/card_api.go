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

// Card maps card to/from json at rest api level
type Card struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

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
	router.Get("/{CARDID}", a.Get)
}

// Get card by id
func (a *CardAPI) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "CARDID")
	op := handlers.Get(id, a, &cardGetMapper{}, a.Logger)
	handlers.Handle(w, r, op)
}

// GetByID implements handlers.GetService
func (a *CardAPI) GetByID(ctx context.Context, id kernel.Id) (interface{}, error) {
	return a.CardService.GetByID(ctx, id)
}

type cardGetMapper struct {
}

// ModelToPayload mapping
func (cardGetMapper) ModelToPayload(m interface{}) interface{} {
	model := m.(*services.CardModel)
	return &Card{
		ID:   string(model.ID),
		Name: model.Name,
	}
}
