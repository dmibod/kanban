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

// CardAPI holds dependencies required by handlers
type CardAPI struct {
	logger.Logger
	services.CardService
}

// CreateCardAPI creates new instance of API
func CreateCardAPI(l logger.Logger, s services.CardService) *CardAPI {
	return &CardAPI{
		Logger:      l,
		CardService: s,
	}
}

// Routes export API router
func (a *CardAPI) Routes(router chi.Router) {
	router.Post("/", a.Create)
	router.Put("/{ID}", a.Update)
	router.Delete("/{ID}", a.Remove)
}

// Create creates new card
func (a *CardAPI) Create(w http.ResponseWriter, r *http.Request) {
	card := &Card{}

	err := mux.JsonRequest(r, card)
	if err != nil {
		a.Errorln("error parsing json", err)
		mux.ErrorResponse(w, http.StatusInternalServerError)
		return
	}

	id, err := a.CardService.CreateCard(r.Context(), &services.CardPayload{Name: card.Name})
	if err != nil {
		a.Errorln("error inserting document", err)
		mux.ErrorResponse(w, http.StatusInternalServerError)
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

	err := mux.JsonRequest(r, card)
	if err != nil {
		a.Errorln("error parsing json", err)
		mux.ErrorResponse(w, http.StatusInternalServerError)
		return
	}

	model, err := a.CardService.UpdateCard(r.Context(), &services.CardModel{ID: kernel.Id(id), Name: card.Name})
	if err != nil {
		a.Errorln("error updating document", err)
		mux.ErrorResponse(w, http.StatusInternalServerError)
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
		a.Errorln("error removing document", err)
		mux.ErrorResponse(w, http.StatusInternalServerError)
		return
	}

	resp := struct {
		ID      string `json:"id"`
		Success bool   `json:"success"`
	}{string(id), true}

	render.JSON(w, r, resp)
}
