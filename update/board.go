package update

import (
	"context"
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

// BoardServiceFactory expected by handler
type BoardServiceFactory interface {
	CreateBoardService(context.Context) services.BoardService
}

// BoardAPI holds dependencies required by handlers
type BoardAPI struct {
	logger.Logger
	BoardServiceFactory
}

// CreateBoardAPI creates API
func CreateBoardAPI(l logger.Logger, f BoardServiceFactory) *BoardAPI {
	return &BoardAPI{
		Logger:              l,
		BoardServiceFactory: f,
	}
}

// Routes export API router
func (a *BoardAPI) Routes(router chi.Router) {
	router.Post("/", a.Create)
	router.Put("/{ID}", a.Update)
	router.Delete("/{ID}", a.Remove)
}

// Create board
func (a *BoardAPI) Create(w http.ResponseWriter, r *http.Request) {
	board := &Board{}

	err := mux.JsonRequest(r, board)
	if err != nil {
		a.Errorln("error parsing json", err)
		mux.ErrorResponse(w, http.StatusBadRequest)
		return
	}

	id, err := a.getService(r).Create(&services.BoardPayload{Name: board.Name})
	if err != nil {
		a.Errorln("error inserting document", err)
		mux.ErrorResponse(w, http.StatusInternalServerError)
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

	err := mux.JsonRequest(r, board)
	if err != nil {
		a.Errorln("error parsing json", err)
		mux.ErrorResponse(w, http.StatusBadRequest)
		return
	}

	model, err := a.getService(r).Update(&services.BoardModel{ID: kernel.Id(id), Name: board.Name})
	if err != nil {
		a.Errorln("error updating document", err)
		mux.ErrorResponse(w, http.StatusInternalServerError)
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

	err := a.getService(r).Remove(kernel.Id(id))
	if err != nil {
		a.Errorln("error removing", err)
		mux.ErrorResponse(w, http.StatusInternalServerError)
		return
	}

	resp := struct {
		ID      string `json:"id"`
		Success bool   `json:"success"`
	}{string(id), true}

	render.JSON(w, r, resp)
}

func (a *BoardAPI) getService(r *http.Request) services.BoardService {
	return a.BoardServiceFactory.CreateBoardService(r.Context())
}
