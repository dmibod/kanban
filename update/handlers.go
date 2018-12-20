package update

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/dmibod/kanban/shared/services"
	"github.com/dmibod/kanban/shared/kernel"
	"net/http"

	"github.com/dmibod/kanban/shared/tools/logger"
	"github.com/dmibod/kanban/shared/tools/mux"
)

// Card maps card to/from json at rest api level
type Card struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// CardService service expected by handler
type CardService interface{
	CreateCard(*services.CardPayload) (kernel.Id, error)
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
	router.Post("/", a.Create)
	return router
}

// Create creates new card
func (a *API) Create(w http.ResponseWriter, r *http.Request) {
	card := &Card{}

	err := mux.JsonRequest(r, card)
	if err != nil {
		a.logger.Errorln("error parsing json", err)
		mux.ErrorResponse(w, http.StatusInternalServerError)		
	}

	id, err := a.service.CreateCard(&services.CardPayload{Name: card.Name})
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
