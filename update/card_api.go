package update

import (
	"context"
	"net/http"

	"github.com/dmibod/kanban/shared/handlers"

	"github.com/dmibod/kanban/shared/kernel"
	"github.com/dmibod/kanban/shared/services"
	"github.com/go-chi/chi"

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
	router.Put("/{CARDID}", a.UpdateCard)
	router.Delete("/{CARDID}", a.RemoveCard)
}

// CreateCard handler
func (a *CardAPI) CreateCard(w http.ResponseWriter, r *http.Request) {
	op := handlers.Create(&Card{}, a, &cardCreateMapper{}, a.Logger)
	handlers.Handle(w, r, op)
}

// UpdateCard handler
func (a *CardAPI) UpdateCard(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "CARDID")
	op := handlers.Update(&Card{ID: id}, a, &cardUpdateMapper{}, a.Logger)
	handlers.Handle(w, r, op)
}

// RemoveCard handler
func (a *CardAPI) RemoveCard(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "CARDID")
	op := handlers.Remove(id, a.CardService, a.Logger)
	handlers.Handle(w, r, op)
}

// Create implements handlers.CreateService
func (a *CardAPI) Create(ctx context.Context, model interface{}) (interface{}, error) {
	return a.CardService.Create(ctx, model.(*services.CardPayload))
}

// Update implements handlers.UpdateService
func (a *CardAPI) Update(ctx context.Context, model interface{}) (interface{}, error) {
	return a.CardService.Update(ctx, model.(*services.CardModel))
}

type cardCreateMapper struct {
}

// PayloadToModel mapping
func (cardCreateMapper) PayloadToModel(p interface{}) interface{} {
	payload := p.(*Card)
	return &services.CardPayload{
		Name: payload.Name,
	}
}

// ModelToPayload mapping
func (cardCreateMapper) ModelToPayload(m interface{}) interface{} {
	model := m.(*services.CardModel)
	return &Card{
		ID:   string(model.ID),
		Name: model.Name,
	}
}

type cardUpdateMapper struct {
}

// PayloadToModel mapping
func (cardUpdateMapper) PayloadToModel(p interface{}) interface{} {
	payload := p.(*Card)
	return &services.CardModel{
		ID:   kernel.ID(payload.ID),
		Name: payload.Name,
	}
}

// ModelToPayload mapping
func (cardUpdateMapper) ModelToPayload(m interface{}) interface{} {
	model := m.(*services.CardModel)
	return &Card{
		ID:   string(model.ID),
		Name: model.Name,
	}
}
