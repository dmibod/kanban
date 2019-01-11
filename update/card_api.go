package update

import (
	"context"
	"net/http"

	"github.com/dmibod/kanban/shared/handlers"

	"github.com/dmibod/kanban/shared/services"
	"github.com/go-chi/chi"

	"github.com/dmibod/kanban/shared/tools/logger"
)

// Card maps card to/from json at rest api level
type Card struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

// CardAPI dependencies
type CardAPI struct {
	services.CardService
	logger.Logger
}

// CreateCardAPI creates API
func CreateCardAPI(s services.CardService, l logger.Logger) *CardAPI {
	return &CardAPI{
		CardService: s,
		Logger:      l,
	}
}

// Routes install handlers
func (a *CardAPI) Routes(router chi.Router) {
	router.Post("/", a.CreateCard)
}

// CreateCard handler
func (a *CardAPI) CreateCard(w http.ResponseWriter, r *http.Request) {
	op := handlers.Create(&Card{}, a, &cardCreateMapper{}, a.Logger)
	handlers.Handle(w, r, op)
}

// Create implements handlers.CreateService
func (a *CardAPI) Create(ctx context.Context, model interface{}) (interface{}, error) {
	return a.CardService.Create(ctx, model.(*services.CardPayload))
}

type cardCreateMapper struct {
}

// PayloadToModel mapping
func (cardCreateMapper) PayloadToModel(p interface{}) interface{} {
	payload := p.(*Card)
	return &services.CardPayload{
		Name:        payload.Name,
		Description: payload.Description,
	}
}

// ModelToPayload mapping
func (cardCreateMapper) ModelToPayload(m interface{}) interface{} {
	model := m.(*services.CardModel)
	return &Card{
		ID:          string(model.ID),
		Name:        model.Name,
		Description: model.Description,
	}
}
