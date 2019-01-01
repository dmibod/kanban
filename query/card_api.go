package query

import (
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
	router.Get("/{ID}", a.Get)
}

// Get card by id
func (a *CardAPI) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "ID")

	model, err := a.CardService.GetCardByID(r.Context(), kernel.Id(id))
	if err != nil {
		a.Errorln("error getting card", err)
		mux.RenderError(w, http.StatusNotFound)
		return
	}

	resp := &Card{
		ID:   string(model.ID),
		Name: model.Name,
	}

	render.JSON(w, r, resp)
}
