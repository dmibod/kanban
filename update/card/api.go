package card

import (
	"github.com/dmibod/kanban/shared/kernel"
	"net/http"

	"github.com/dmibod/kanban/shared/handlers"

	"github.com/dmibod/kanban/shared/services/card"
	"github.com/go-chi/chi"

	"github.com/dmibod/kanban/shared/tools/logger"
)

// API dependencies
type API struct {
	card.Service
	logger.Logger
}

// CreateAPI creates API
func CreateAPI(s card.Service, l logger.Logger) *API {
	return &API{
		Service: s,
		Logger:  l,
	}
}

// Routes install handlers
func (a *API) Routes(router chi.Router) {
	router.Post("/{BOARDID}/cards", a.CreateCard)
}

// CreateCard handler
func (a *API) CreateCard(w http.ResponseWriter, r *http.Request) {
	op := handlers.Create(&Card{}, a, &cardCreateMapper{}, a.Logger)
	handlers.Handle(w, r, op)
}

// Create implements handlers.CreateService
func (a *API) Create(r *http.Request, model interface{}) (interface{}, error) {
	boardID := chi.URLParam(r, "BOARDID")
	cardID, err := a.Service.Create(r.Context(), kernel.ID(boardID), model.(*card.CreateModel))
	if err != nil {
		return nil, err
	}
	return a.Service.GetByID(r.Context(), cardID.WithSet(kernel.ID(boardID)))
}
