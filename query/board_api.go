package query

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
	router.Get("/{ID}", a.Get)
	router.Get("/", a.All)
}

// Get by id
func (a *BoardAPI) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "ID")

	model, err := a.BoardService.GetByID(r.Context(), kernel.Id(id))
	if err != nil {
		a.Errorln(err)
		mux.RenderError(w, http.StatusNotFound)
		return
	}

	resp := &Board{
		ID:   string(model.ID),
		Name: model.Name,
	}

	render.JSON(w, r, resp)
}

// All boards
func (a *BoardAPI) All(w http.ResponseWriter, r *http.Request) {
	models, err := a.BoardService.GetAll(r.Context())
	if err != nil {
		a.Errorln(err)
		mux.RenderError(w, http.StatusInternalServerError)
		return
	}

	resp := []*Board{}
	for _, model := range models {
		board := &Board{
			ID:   string(model.ID),
			Name: model.Name,
		}

		resp = append(resp, board)
	}

	render.JSON(w, r, resp)
}
