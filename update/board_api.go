package update

import (
	"github.com/dmibod/kanban/shared/kernel"
	"github.com/dmibod/kanban/shared/services"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"net/http"

	"github.com/dmibod/kanban/shared/tools/logger"
	"github.com/dmibod/kanban/shared/tools/mux"
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
	router.Post("/", a.Create)
	router.Put("/{ID}", a.Update)
	router.Delete("/{ID}", a.Remove)
}

// Create board
func (a *BoardAPI) Create(w http.ResponseWriter, r *http.Request) {
	board := &Board{}

	err := mux.ParseJSON(r, board)
	if err != nil {
		a.Errorln(err)
		mux.RenderError(w, http.StatusBadRequest)
		return
	}

	id, err := a.BoardService.Create(r.Context(), &services.BoardPayload{Name: board.Name})
	if err != nil {
		a.Errorln(err)
		mux.RenderError(w, http.StatusInternalServerError)
		return
	}

	resp := struct {
		ID      string `json:"id"`
		Success bool   `json:"success"`
	}{string(id), true}

	render.JSON(w, r, resp)
}

// Update board
func (a *BoardAPI) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "ID")
	board := &Board{}

	err := mux.ParseJSON(r, board)
	if err != nil {
		a.Errorln(err)
		mux.RenderError(w, http.StatusBadRequest)
		return
	}

	model, err := a.BoardService.Update(r.Context(), &services.BoardModel{ID: kernel.Id(id), Name: board.Name})
	if err != nil {
		a.Errorln(err)
		mux.RenderError(w, http.StatusInternalServerError)
		return
	}

	resp := &Board{
		ID:   string(model.ID),
		Name: model.Name,
	}

	render.JSON(w, r, resp)
}

// Remove board
func (a *BoardAPI) Remove(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "ID")

	err := a.BoardService.Remove(r.Context(), kernel.Id(id))
	if err != nil {
		a.Errorln(err)
		mux.RenderError(w, http.StatusInternalServerError)
		return
	}

	resp := struct {
		ID      string `json:"id"`
		Success bool   `json:"success"`
	}{string(id), true}

	render.JSON(w, r, resp)
}
