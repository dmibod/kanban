package query

import (
	"context"
	"net/http"

	"github.com/dmibod/kanban/shared/handlers"
	"github.com/dmibod/kanban/shared/kernel"
	"github.com/dmibod/kanban/shared/services"
	"github.com/go-chi/chi"

	"github.com/dmibod/kanban/shared/tools/logger"
)

// Board api
type Board struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// BoardAPI dependencies
type BoardAPI struct {
	services.BoardService
	logger.Logger
}

// CreateBoardAPI creates API
func CreateBoardAPI(s services.BoardService, l logger.Logger) *BoardAPI {
	return &BoardAPI{
		BoardService: s,
		Logger:       l,
	}
}

// Routes install handlers
func (a *BoardAPI) Routes(router chi.Router) {
	router.Get("/{BOARDID}", a.Get)
	router.Get("/", a.All)
}

// Get by id
func (a *BoardAPI) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "BOARDID")
	op := handlers.Get(id, a, &boardGetMapper{}, a.Logger)
	handlers.Handle(w, r, op)
}

// All boards
func (a *BoardAPI) All(w http.ResponseWriter, r *http.Request) {
	op := handlers.All(a, &boardGetMapper{}, a.Logger)
	handlers.Handle(w, r, op)
}

// GetAll implements handlers.AllService
func (a *BoardAPI) GetAll(ctx context.Context) ([]interface{}, error) {
	models, err := a.BoardService.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	items := []interface{}{}
	for _, model := range models {
		items = append(items, model)
	}
	return items, nil
}

// GetByID implements handlers.GetService
func (a *BoardAPI) GetByID(ctx context.Context, id kernel.Id) (interface{}, error) {
	return a.BoardService.GetByID(ctx, id)
}

type boardGetMapper struct {
}

// ModelToPayload mapping
func (boardGetMapper) ModelToPayload(m interface{}) interface{} {
	model := m.(*services.BoardModel)
	return &Board{
		ID:   string(model.ID),
		Name: model.Name,
	}
}
