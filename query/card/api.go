package card

import (
	"net/http"

	"github.com/dmibod/kanban/shared/handlers"

	"github.com/dmibod/kanban/shared/services/card"
	"github.com/go-chi/chi"

	"github.com/dmibod/kanban/shared/kernel"
	"github.com/dmibod/kanban/shared/tools/logger"
)

// API dependencies
type API struct {
	card.Service
	logger.Logger
}

// CreateAPI creates new instance of API
func CreateAPI(s card.Service, l logger.Logger) *API {
	return &API{
		Service: s,
		Logger:  l,
	}
}

// Routes install API handlers
func (a *API) Routes(router chi.Router) {
	router.Get("/{BOARDID}/lanes/{LANEID}/cards", a.List)
	router.Get("/{BOARDID}/cards/{CARDID}", a.Get)
}

// List cards by lane
func (a *API) List(w http.ResponseWriter, r *http.Request) {
	op := handlers.List(a, &CardGetMapper{}, a.Logger)
	handlers.Handle(w, r, op)
}

// Get card by id
func (a *API) Get(w http.ResponseWriter, r *http.Request) {
	op := handlers.Get(a, &CardGetMapper{}, a.Logger)
	handlers.Handle(w, r, op)
}

// GetList implements handlers.ListService
func (a *API) GetList(r *http.Request) ([]interface{}, error) {
	boardID := chi.URLParam(r, "BOARDID")
	laneID := chi.URLParam(r, "LANEID")
	if models, err := a.Service.GetByLaneID(r.Context(), kernel.ID(boardID).WithID(kernel.ID(laneID))); err != nil {
		return nil, err
	} else {
		mapper := CardGetMapper{}
		return mapper.List(models), nil
	}
}

// GetOne implements handlers.GetService
func (a *API) GetOne(r *http.Request) (interface{}, error) {
	boardID := chi.URLParam(r, "BOARDID")
	cardID := chi.URLParam(r, "CARDID")
	return a.Service.GetByID(r.Context(), kernel.ID(boardID).WithID(kernel.ID(cardID)))
}
