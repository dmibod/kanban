package query

import (
	"context"
	"net/http"

	"github.com/dmibod/kanban/shared/services"
	"github.com/dmibod/kanban/shared/tools/mux"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"

	"github.com/dmibod/kanban/shared/kernel"
	"github.com/dmibod/kanban/shared/tools/logger"
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
	router.Get("/{ID}", a.Get)
	router.Get("/", a.All)
}

// Get - gets card by id
func (a *API) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "ID")

	model, err := a.getService(r.Context()).GetCardByID(kernel.Id(id))
	if err != nil {
		a.logger.Errorln("error getting card", err)
		mux.ErrorResponse(w, http.StatusInternalServerError)
		return
	}

	resp := &Card{
		ID:   string(model.ID),
		Name: model.Name,
	}

	render.JSON(w, r, resp)
}

// All - gets all cards
func (a *API) All(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "ID")

	model, err := a.getService(r.Context()).GetCardByID(kernel.Id(id))
	if err != nil {
		a.logger.Errorln("error getting card", err)
		mux.ErrorResponse(w, http.StatusInternalServerError)
		return
	}

	resp := &Card{
		ID:   string(model.ID),
		Name: model.Name,
	}

	render.JSON(w, r, resp)
}

func (a *API) getService(c context.Context) services.CardService {
	return a.factory.CreateCardService(c)
}
