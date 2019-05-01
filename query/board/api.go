package board

import (
	"net/http"

	"github.com/dmibod/kanban/shared/handlers"
	"github.com/dmibod/kanban/shared/kernel"
	"github.com/dmibod/kanban/shared/services/board"
	"github.com/go-chi/chi"

	"github.com/dmibod/kanban/shared/tools/logger"
)

// API dependencies
type API struct {
	board.Service
	logger.Logger
}

// CreateAPI creates API
func CreateAPI(s board.Service, l logger.Logger) *API {
	return &API{
		Service: s,
		Logger:  l,
	}
}

// Routes install handlers
func (a *API) Routes(router chi.Router) {
	router.Get("/", a.List)
	router.Get("/{BOARDID}", a.Get)
}

// List boards
func (a *API) List(w http.ResponseWriter, r *http.Request) {
	op := handlers.List(a, &ListModelMapper{}, a.Logger)
	handlers.Handle(w, r, op)
}

// GetList for board
func (a *API) GetList(r *http.Request) ([]interface{}, error) {
	owner := r.URL.Query().Get("owner")
	if models, err := a.Service.GetByOwner(r.Context(), owner); err != nil {
		return nil, err
	} else {
		mapper := &ListModelMapper{}
		return mapper.List(models), nil
	}
}

// Get by id
func (a *API) Get(w http.ResponseWriter, r *http.Request) {
	op := handlers.Get(a, &ListModelMapper{}, a.Logger)
	handlers.Handle(w, r, op)
}

// GetOne implements handlers.GetService
func (a *API) GetOne(r *http.Request) (interface{}, error) {
	id := chi.URLParam(r, "BOARDID")
	return a.Service.GetByID(r.Context(), kernel.ID(id))
}
