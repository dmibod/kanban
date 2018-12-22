package update

import (
	"github.com/dmibod/kanban/shared/kernel"
	"context"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/dmibod/kanban/shared/services"
	"net/http"

	"github.com/dmibod/kanban/shared/tools/logger"
	"github.com/dmibod/kanban/shared/tools/mux"
)

// Card maps card to/from json at rest api level
type Card struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// ServiceFactory factory expected by handler
type ServiceFactory interface {
	CreateCardService(context.Context) services.CardService
}

// API holds dependencies required by handlers
type API struct {
	logger  logger.Logger
	factory ServiceFactory
}

// CreateAPI creates new instance of API
func CreateAPI(l logger.Logger, f ServiceFactory) *API {
	return &API{
		logger:  l,
		factory: f,
	}
}

// Routes export API router
func (a *API) Routes(router *chi.Mux) {
	router.Post("/", a.Create)
	router.Put("/{ID}", a.Update)
	router.Delete("/{ID}", a.Remove)
}

// Create creates new card
func (a *API) Create(w http.ResponseWriter, r *http.Request) {
	card := &Card{}

	err := mux.JsonRequest(r, card)
	if err != nil {
		a.logger.Errorln("error parsing json", err)
		mux.ErrorResponse(w, http.StatusInternalServerError)
		return		
	}

	id, err := a.getService(r.Context()).CreateCard(&services.CardPayload{Name: card.Name})
	if err != nil {
		a.logger.Errorln("error inserting document", err)
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
func (a *API) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "ID")
	card := &Card{}

	err := mux.JsonRequest(r, card)
	if err != nil {
		a.logger.Errorln("error parsing json", err)
		mux.ErrorResponse(w, http.StatusInternalServerError)
		return		
	}

	model, err := a.getService(r.Context()).UpdateCard(&services.CardModel{ID: kernel.Id(id), Name: card.Name})
	if err != nil {
		a.logger.Errorln("error updating document", err)
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
func (a *API) Remove(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "ID")

	err := a.getService(r.Context()).RemoveCard(kernel.Id(id))
	if err != nil {
		a.logger.Errorln("error removing document", err)
		mux.ErrorResponse(w, http.StatusInternalServerError)
		return
	}

	resp := struct {
		ID      string `json:"id"`
		Success bool   `json:"success"`
	}{string(id), true}

	render.JSON(w, r, resp)
}

func (a *API) getService(c context.Context) services.CardService {
	return a.factory.CreateCardService(c)
}
