package update

import (
	"context"
	"net/http"

	"github.com/dmibod/kanban/shared/tools/mux"
	"github.com/go-chi/render"

	"github.com/dmibod/kanban/shared/handlers"
	"github.com/dmibod/kanban/shared/kernel"
	"github.com/dmibod/kanban/shared/services"
	"github.com/go-chi/chi"

	"github.com/dmibod/kanban/shared/tools/logger"
)

// Board api
type Board struct {
	ID     string `json:"id,omitempty"`
	Name   string `json:"name,omitempty"`
	Layout string `json:"layout,omitempty"`
	Owner  string `json:"owner,omitempty"`
	Shared bool   `json:"shared,omitempty"`
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
	router.Put("/{BOARDID}/rename", a.RenameBoard)
	router.Put("/{BOARDID}/share", a.ShareBoard)
	router.Delete("/{BOARDID}", a.RemoveBoard)
}

// CreateBoard handler
func (a *BoardAPI) CreateBoard(w http.ResponseWriter, r *http.Request) {
	op := handlers.Create(&Board{}, a, &boardMapper{}, a.Logger)
	handlers.Handle(w, r, op)
}

// RenameBoard handler
func (a *BoardAPI) RenameBoard(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "BOARDID")
	payload := &struct {
		Name string `json:"name,omitempty"`
	}{}
	if err := mux.ParseJSON(r, payload); err != nil {
		a.Errorln(err)
		mux.RenderError(w, http.StatusBadRequest)
		return
	}
	model, err := a.BoardService.Rename(r.Context(), kernel.Id(id), payload.Name)
	if err != nil {
		a.Errorln(err)
		mux.RenderError(w, http.StatusInternalServerError)
		return
	}

	mapper := &boardMapper{}
	render.JSON(w, r, mapper.ModelToPayload(model))
}

// ShareBoard handler
func (a *BoardAPI) ShareBoard(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "BOARDID")
	payload := &struct {
		Shared bool `json:"shared,omitempty"`
	}{}
	if err := mux.ParseJSON(r, payload); err != nil {
		a.Errorln(err)
		mux.RenderError(w, http.StatusBadRequest)
		return
	}

	model, err := a.BoardService.Share(r.Context(), kernel.Id(id), payload.Shared)
	if err != nil {
		a.Errorln(err)
		mux.RenderError(w, http.StatusInternalServerError)
		return
	}
	mapper := &boardMapper{}
	render.JSON(w, r, mapper.ModelToPayload(model))
}

// RemoveBoard handler
func (a *BoardAPI) RemoveBoard(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "BOARDID")
	op := handlers.Remove(id, a.BoardService, a.Logger)
	handlers.Handle(w, r, op)
}

// Create implements handlers.CreateService
func (a *BoardAPI) Create(ctx context.Context, model interface{}) (interface{}, error) {
	return a.BoardService.Create(ctx, model.(*services.BoardPayload))
}

type boardMapper struct {
}

// PayloadToModel mapping
func (boardMapper) PayloadToModel(p interface{}) interface{} {
	payload := p.(*Board)

	return &services.BoardPayload{
		Name:   payload.Name,
		Layout: payload.Layout,
		Owner:  payload.Owner,
	}
}

func (boardMapper) ModelToPayload(m interface{}) interface{} {
	model := m.(*services.BoardModel)

	return &Board{
		ID:     string(model.ID),
		Name:   model.Name,
		Layout: model.Layout,
		Owner:  model.Owner,
		Shared: model.Shared,
	}
}
