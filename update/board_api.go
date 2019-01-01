package update

import (
	"context"
	"github.com/dmibod/kanban/shared/handlers"
	"github.com/dmibod/kanban/shared/kernel"
	"github.com/dmibod/kanban/shared/services"
	"github.com/go-chi/chi"
	"net/http"

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
	router.Post("/", a.CreateBoard)
	router.Put("/{BOARDID}", a.UpdateBoard)
	router.Delete("/{BOARDID}", a.RemoveBoard)
}

// CreateBoard handler
func (a *BoardAPI) CreateBoard(w http.ResponseWriter, r *http.Request) {
	op := handlers.Create(&Board{}, a, &cardCreateMapper{}, a.Logger)
	handlers.Handle(w, r, op)
}

// UpdateBoard handler
func (a *BoardAPI) UpdateBoard(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "BOARDID")
	op := handlers.Update(&Board{ID: id}, a, &boardUpdateMapper{}, a.Logger)
	handlers.Handle(w, r, op)
}

// RemoveBoard handler
func (a *BoardAPI) RemoveBoard(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "BOARDID")
	op := handlers.Remove(id, a.BoardService, a.Logger)
	handlers.Handle(w, r, op)
}

// Create implements handlers.CreateService
func (a *BoardAPI) Create(ctx context.Context, model interface{}) (kernel.Id, error) {
	return a.BoardService.Create(ctx, model.(*services.BoardPayload))
}

// Update implements handlers.UpdateService
func (a *BoardAPI) Update(ctx context.Context, model interface{}) (interface{}, error) {
	return a.BoardService.Update(ctx, model.(*services.BoardModel))
}

type boardCreateMapper struct {
}

// PayloadToModel mapping
func (boardCreateMapper) PayloadToModel(p interface{}) interface{} {
	payload := p.(*Board)
	return &services.BoardPayload{
		Name: payload.Name,
	}
}

type boardUpdateMapper struct {
}

// PayloadToModel mapping
func (boardUpdateMapper) PayloadToModel(p interface{}) interface{} {
	payload := p.(*Board)
	return &services.BoardModel{
		ID:   kernel.Id(payload.ID),
		Name: payload.Name,
	}
}

// ModelToPayload mapping
func (boardUpdateMapper) ModelToPayload(m interface{}) interface{} {
	model := m.(*services.BoardModel)
	return &Board{
		ID:   string(model.ID),
		Name: model.Name,
	}
}
