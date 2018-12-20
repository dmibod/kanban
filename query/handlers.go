package query

import (
	"github.com/go-chi/render"
	"github.com/go-chi/chi"
	"github.com/dmibod/kanban/shared/services"
	"net/http"

	"github.com/dmibod/kanban/shared/kernel"
	"github.com/dmibod/kanban/shared/tools/logger"
)

// Card maps card to/from json at rest api level
type Card struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// CardService service expected by handler
type CardService interface {
	GetCardByID(kernel.Id) (*services.CardModel, error)
}

// API holds dependencies required by handlers
type API struct {
	logger  logger.Logger
	service CardService
}

// CreateAPI creates new instance of API
func CreateAPI(l logger.Logger, s CardService) *API {
	return &API{
		logger:  l,
		service: s,
	}
}

// Routes export API router
func (a *API) Routes() *chi.Mux {
	router := chi.NewRouter()
	router.Get("/{ID}", a.Get)
	router.Get("/", a.All)
	return router
}

// Get - gets card by id
func (a *API) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "ID")

	model, err := a.service.GetCardByID(kernel.Id(id))
	if err != nil {
		a.logger.Errorln("error getting card", err)
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

	model, err := a.service.GetCardByID(kernel.Id(id))
	if err != nil {
		a.logger.Errorln("error getting card", err)
	}

	resp := &Card{
		ID:   string(model.ID),
		Name: model.Name,
	}

	render.JSON(w, r, resp)
}
