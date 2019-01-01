package update

import (
	"net/http"

	"github.com/dmibod/kanban/shared/kernel"
	"github.com/dmibod/kanban/shared/services"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"

	"github.com/dmibod/kanban/shared/tools/logger"
	"github.com/dmibod/kanban/shared/tools/mux"
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
	router.Post("/", a.Create)
	router.Put("/{ID}", a.Update)
	router.Delete("/{ID}", a.Remove)
}

// Create creates new card
func (a *CardAPI) Create(w http.ResponseWriter, r *http.Request) {
	card := &Card{}

	err := mux.ParseJSON(r, card)
	if err != nil {
		a.Errorln(err)
		mux.RenderError(w, http.StatusInternalServerError)
		return
	}

	id, err := a.CardService.CreateCard(r.Context(), &services.CardPayload{Name: card.Name})
	if err != nil {
		a.Errorln(err)
		mux.RenderError(w, http.StatusInternalServerError)
		return
	}

	resp := struct {
		ID      string `json:"id"`
		Success bool   `json:"success"`
	}{string(id), true}

	render.JSON(w, r, resp)
}

// Update updates card
func (a *CardAPI) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "ID")
	card := &Card{}

	err := mux.ParseJSON(r, card)
	if err != nil {
		a.Errorln(err)
		mux.RenderError(w, http.StatusInternalServerError)
		return
	}

	model, err := a.CardService.UpdateCard(r.Context(), &services.CardModel{ID: kernel.Id(id), Name: card.Name})
	if err != nil {
		a.Errorln(err)
		mux.RenderError(w, http.StatusInternalServerError)
		return
	}

	resp := &Card{
		ID:   string(model.ID),
		Name: model.Name,
	}

	render.JSON(w, r, resp)
}

// Remove removes card
func (a *CardAPI) Remove(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "ID")

	err := a.CardService.RemoveCard(r.Context(), kernel.Id(id))
	if err != nil {
		a.Errorln(err)
		mux.RenderError(w, http.StatusInternalServerError)
		return
	}

	resp := struct {
		ID      string `json:"id"`
		Success bool   `json:"success"`
	}{string(id), true}

	render.JSON(w, r, resp)
}
