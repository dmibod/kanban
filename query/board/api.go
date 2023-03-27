package board

import (
	"github.com/dmibod/kanban/shared/handlers"
	"github.com/dmibod/kanban/shared/kernel"
	"github.com/dmibod/kanban/shared/services/board"
	"github.com/dmibod/kanban/shared/tools/logger"
	"github.com/go-chi/chi"
	"net/http"
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
	router.Get("/{ID}", a.Get)
	/*
		keyFunc := func(r *http.Request) uint64 {
			boardId := chi.URLParam(r, "ID")
			return stampede.StringToHash(strings.ToLower(boardId))
		}

		cached := stampede.HandlerWithKey(10*1024, 30*time.Second, keyFunc)

		router.With(cached).Get("/{ID}", a.Get)
	*/
}

// List boards
func (a *API) List(w http.ResponseWriter, r *http.Request) {
	op := handlers.List(a, &ListModelMapper{}, a.Logger)
	handlers.Handle(w, r, op)
}

// GetList implements handlers.ListService
func (a *API) GetList(r *http.Request) ([]interface{}, error) {
	owner := r.URL.Query().Get("owner")
	if models, err := a.Service.GetByOwner(r.Context(), owner); err != nil {
		return nil, err
	} else {
		mapper := &ListModelMapper{}
		return mapper.List(models), nil
	}
}

// Get board by id
func (a *API) Get(w http.ResponseWriter, r *http.Request) {
	op := handlers.Get(a, &ModelMapper{}, a.Logger)
	handlers.Handle(w, r, op)
}

// GetOne implements handlers.GetService
func (a *API) GetOne(r *http.Request) (interface{}, error) {
	id := chi.URLParam(r, "ID")

	a.Logger.Debugf("get board %v", id)

	return a.Service.GetByID(r.Context(), kernel.ID(id))
}
