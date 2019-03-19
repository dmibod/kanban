package board

import (
	"context"
	"net/http"

	"github.com/dmibod/kanban/shared/tools/mux"
	"github.com/go-chi/render"

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
	router.Post("/", a.CreateBoard)
	router.Put("/{BOARDID}/rename", a.RenameBoard)
	router.Put("/{BOARDID}/share", a.ShareBoard)
	router.Delete("/{BOARDID}", a.RemoveBoard)
}

// CreateBoard handler
func (a *API) CreateBoard(w http.ResponseWriter, r *http.Request) {
	op := handlers.Create(&Board{}, a, &boardMapper{}, a.Logger)
	handlers.Handle(w, r, op)
}

// RenameBoard handler
func (a *API) RenameBoard(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "BOARDID")
	payload := &struct {
		Name string `json:"name,omitempty"`
	}{}
	if err := mux.ParseJSON(r, payload); err != nil {
		a.Errorln(err)
		mux.RenderError(w, http.StatusBadRequest)
		return
	}
	
	if err := a.Service.Name(r.Context(), kernel.ID(id), payload.Name); err != nil {
		a.Errorln(err)
		mux.RenderError(w, http.StatusInternalServerError)
		return
	}

	model, err := a.Service.GetByID(r.Context(), kernel.ID(id))
	if err != nil {
		a.Errorln(err)
		mux.RenderError(w, http.StatusInternalServerError)
		return
	}

	mapper := &boardMapper{}
	render.JSON(w, r, mapper.ModelToPayload(model))
}

// ShareBoard handler
func (a *API) ShareBoard(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "BOARDID")
	payload := &struct {
		Shared bool `json:"shared,omitempty"`
	}{}
	if err := mux.ParseJSON(r, payload); err != nil {
		a.Errorln(err)
		mux.RenderError(w, http.StatusBadRequest)
		return
	}

	if err := a.Service.Share(r.Context(), kernel.ID(id), payload.Shared); err != nil {
		a.Errorln(err)
		mux.RenderError(w, http.StatusInternalServerError)
		return
	}

	model, err := a.Service.GetByID(r.Context(), kernel.ID(id))
	if err != nil {
		a.Errorln(err)
		mux.RenderError(w, http.StatusInternalServerError)
		return
	}

	mapper := &boardMapper{}
	render.JSON(w, r, mapper.ModelToPayload(model))
}

// RemoveBoard handler
func (a *API) RemoveBoard(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "BOARDID")
	op := handlers.Remove(id, a.Service, a.Logger)
	handlers.Handle(w, r, op)
}

// Create implements handlers.CreateService
func (a *API) Create(ctx context.Context, model interface{}) (interface{}, error) {
	return a.Service.Create(ctx, model.(*board.CreateModel))
}
